package role

import (
	"fmt"
	"net/http"

	"github.com/Nasa28/CommerceCore/types"
	"github.com/Nasa28/CommerceCore/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type RoleHandler struct {
	store types.RoleStore
}

func NewRoleHandler(store types.RoleStore) *RoleHandler {
	return &RoleHandler{store: store}
}

func (h *RoleHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/roles", h.handleCreateRole).Methods("POST")
	router.HandleFunc("/roles", h.handleGetRoles).Methods("GET")
}

func (h *RoleHandler) handleCreateRole(w http.ResponseWriter, r *http.Request) {

	var payload types.Role
	// Parse the request body
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := h.store.CreateRole(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"Message": "Role Created succesfully"})
}
func (h *RoleHandler) handleGetRoles(w http.ResponseWriter, _ *http.Request) {
	// Fetch all roles from the store
	roles, err := h.store.GetAllRoles()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with the retrieved roles
	utils.WriteJSON(w, http.StatusOK, roles)
}
