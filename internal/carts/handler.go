package carts

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/inflame-ue/gocommerce/internal/auth"
	"github.com/inflame-ue/gocommerce/internal/response"
)

func (ch *CartHandler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := userClaims["sub"].(int)
	
	items, err := ch.ListCartItems(r.Context(), userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch cart items from the database"})
		return
	}

	resp := struct {
		CartItems []CartItem `json:"cart_items"`
	}{
		CartItems: items,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

func (ch *CartHandler) HandleAddProductToCart(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := userClaims["sub"].(int)


	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the productID is invalid"})
		return
	}

	err = ch.AddCartItem(r.Context(), userID, productID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not add the product to cart, something went wrong"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "product was successfully added to the cart"})
}

func (ch *CartHandler) HandleDeleteProductFromCart(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := userClaims["sub"].(int)

	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the productID is invalid"})
		return
	}

	err = ch.DeleteCartItem(r.Context(), userID, productID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete the product from cart, something went wrong"})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "product was successfully deleted from the cart"})
}
