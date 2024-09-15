package handler

import (
	"be-golang-chapter-24/RESTFull-API/collections"
	"be-golang-chapter-24/RESTFull-API/service"
	"be-golang-chapter-24/RESTFull-API/validation"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := collections.Response{
		Status:  "success",
		Message: "Products retrieved successfully",
		Data:    products,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response := collections.Response{
			Status:  "error",
			Message: "Invalid product ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	product, err := h.service.GetProductByID(id)
	if err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := collections.Response{
		Status:  "success",
		Message: "Product retrieved successfully",
		Data:    product,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product collections.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: "Invalid request payload",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validation.ValidateProduct(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.CreateProduct(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := collections.Response{
		Status:  "success",
		Message: "Product created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response := collections.Response{
			Status:  "error",
			Message: "Invalid product ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	var product collections.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: "Invalid request payload",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	product.ID = id

	if err := validation.ValidateProduct(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.UpdateProduct(&product); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := collections.Response{
		Status:  "success",
		Message: "Product updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response := collections.Response{
			Status:  "error",
			Message: "Invalid product ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		response := collections.Response{
			Status:  "error",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := collections.Response{
		Status:  "success",
		Message: "Product deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
