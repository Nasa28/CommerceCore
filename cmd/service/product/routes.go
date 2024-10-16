package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nasa28/CommerceCore/types"
	"github.com/Nasa28/CommerceCore/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repository types.ProductRepository
}

func NewProductHandler(repository types.ProductRepository) *ProductHandler {
	return &ProductHandler{repository: repository}
}

func (p ProductHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", p.handleCreateproduct).Methods("POST")
	router.HandleFunc("/products/{id}", p.handleGetProductByID).Methods("GET")
	router.HandleFunc("/products/{id}", p.handleProductUpdate).Methods("PATCH")
}

func (p *ProductHandler) handleCreateproduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload

	// Parse the request body
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := p.repository.CreateProduct(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"Message": "Product added succesfully"})
}

func (p *ProductHandler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	// Convert the idStr to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Fetch the product by ID
	product, err := p.repository.GetProductByID(id)
	if err != nil {
		// Handle error if product is not found
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Create a flattened response

	// Send the product details as a JSON response
	utils.WriteJSON(w, http.StatusOK, product)
}

func (p *ProductHandler) handleProductUpdate(w http.ResponseWriter, r *http.Request) {

	var payload types.ProductAndInventoryUpdate

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := p.repository.UpdateProduct(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
