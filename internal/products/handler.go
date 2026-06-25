package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/inflame-ue/gocommerce/internal/response"
	"github.com/jackc/pgx/v5"
)

func (ph *ProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.ListProducts(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve products from the database"})
		return
	}

	resp := struct {
		Products []productModel `json:"products"`
	}{
		Products: products,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

func (ph *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "productID parameter must be an integer"})
		return
	}

	product, err := ph.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "no product with such an ID exists"})
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch the product record"})
		return
	}

	response.WriteJSON(w, http.StatusOK, product)
}

func (ph *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var productReq productModel
	if err := json.NewDecoder(r.Body).Decode(&productReq); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the provided JSON could not be parsed"})
		return
	}

	productID, err := ph.CreateProduct(r.Context(), productReq.Name, productReq.Price, productReq.Stock)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create the product record"})
		return
	}

	product, err := ph.GetProductByID(r.Context(), productID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch the created product by ID"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, product)
}

func (ph *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var productReq productModel
	if err := json.NewDecoder(r.Body).Decode(&productReq); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the provided JSON could not be parsed"})
		return
	}

	product, err := ph.UpdateProductByID(r.Context(), productReq.ID, productReq.Name, productReq.Price, productReq.Stock)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not updated the product record"})
		return
	}

	response.WriteJSON(w, http.StatusOK, product)
}

func (ph *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "productID parameter must be an integer"})
		return
	}

	affected, err := ph.DeleteProductByID(r.Context(), productID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete the record from the database"})
		return
	}
	if affected == 0 {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "no product record with the provided productID is available"})
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("product record with id %d was successfully deleted", productID)})
}
