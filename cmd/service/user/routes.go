package user

import (
	"fmt"
	"net/http"

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

	err = h.store.CreateUser(types.User{
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
