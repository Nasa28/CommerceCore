package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nasa28/CommerceCore/cmd/service/auth"
	"github.com/Nasa28/CommerceCore/config"
	"github.com/Nasa28/CommerceCore/types"
	"github.com/Nasa28/CommerceCore/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/users/{id}", h.handleGetUserByID).Methods("GET")
	router.HandleFunc("/users", h.handleGetUsers).Methods("GET")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}
	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.RegisterUserPayload{
		FirstName: payload.FirstName,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
		State:     payload.State,
		Password:  hashedPassword,
		Country:   payload.Country,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	secret := []byte(config.Env.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, fmt.Errorf("user not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword([]byte(payload.Password), []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}
	secret := []byte(config.Env.JWTSecret)
	token, err := auth.CreateJWT(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
func (h *Handler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	// Convert the idStr to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Fetch the product by ID
	product, err := h.store.GetUserByID(id)
	if err != nil {
		// Handle error if product is not found
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Create a flattened response

	// Send the product details as a JSON response
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offsetStr := query.Get("offset")
	offset, err := strconv.Atoi(offsetStr)

	if err != nil || offset < 0 {
		offset = 0
	}

	limitStr := query.Get("limit")

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit < 0 {
		limit = 20
	}

	users, err := h.store.ListUsers(limit, offset)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
