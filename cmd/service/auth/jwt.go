package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Nasa28/CommerceCore/config"
	"github.com/Nasa28/CommerceCore/types"
	"github.com/Nasa28/CommerceCore/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "UserID"

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Env.JWTTokenExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ProtectedRoute(handleFunc http.HandlerFunc, store types.UserStore, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT Token from request
		tokenString := extractToken(r)

		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			AccessDenied(w)
			return
		}

		// Token is valid, extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract expiry and check if it's valid
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					log.Println("Token has expired")
					AccessDenied(w)
					return
				}
			} else {
				log.Println("Token does not contain an expiry claim")
				AccessDenied(w)
				return
			}

			// Extract userID from claims
			userIDStr, ok := claims["userId"].(string)
			if !ok {
				log.Println("userId claim is missing or invalid")
				AccessDenied(w)
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Printf("Failed to convert userId: %v", err)
				AccessDenied(w)
				return
			}

			if store == nil {
				log.Println("User store is nil")
				AccessDenied(w)
				return
			}
			// Fetch user from the repository
			user, err := store.GetUserByID(userID)
			if err != nil {
				log.Printf("Failed to get user by ID: %v", err)
				AccessDenied(w)
				return
			}

			// Check if the user's role is in the allowedRoles list
			if !roleIsAllowed(user.Role, allowedRoles) {
				log.Printf("Access denied: User role '%s' is not allowed", user.Role)

				w.WriteHeader(http.StatusForbidden)

				response := map[string]string{"Message": "Forbidden: You are not authorized to access this resource"}

				// Encode the map as JSON and write it to the response
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				return
			}
			//Set user in context
			ctx := context.WithValue(r.Context(), UserKey, user)
			r = r.WithContext(ctx)

			// // Call the next handler with the updated context
			handleFunc(w, r)
		} else {
			log.Println("Token is invalid or does not contain claims")
			AccessDenied(w)
			return
		}
	}
}

func extractToken(r *http.Request) string {
	authToken := r.Header.Get("Authorization")
	if authToken != "" && strings.HasPrefix(authToken, "Bearer ") {
		return strings.TrimPrefix(authToken, "Bearer ")
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Env.JWTSecret), nil
	})
}

func AccessDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("access denied"))
}

func GetUserFromContext(ctx context.Context) *types.User {
	user, ok := ctx.Value(UserKey).(*types.User)
	if !ok {
		return nil
	}
	return user
}

func roleIsAllowed(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
