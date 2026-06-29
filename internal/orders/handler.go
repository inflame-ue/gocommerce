package orders

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/inflame-ue/gocommerce/internal/auth"
	"github.com/inflame-ue/gocommerce/internal/response"
	"github.com/jackc/pgx/v5"
)

func (oh *OrderHandler) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := int(userClaims["sub"].(float64))

	orders, err := oh.ListOrders(r.Context(), userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not retrieve the orders associated with the current user"})
		return
	}

	resp := struct {
		Orders []OrderModel `json:"orders"`
	}{
		Orders: orders,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

func (oh *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := int(userClaims["sub"].(float64))
	
	orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "orderID is in wrong format"})
		return
	}

	order, err := oh.GetOrderByID(r.Context(), userID, orderID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not retrieve the orders associated with the current user"})
		return
	}

	response.WriteJSON(w, http.StatusOK, order)
}

func (oh *OrderHandler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := int(userClaims["sub"].(float64))
	isAdmin := userClaims["is_admin"].(bool)

	if !isAdmin {
		response.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "you must have admin rights to access this endpoint"})
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "orderID is in wrong format"})
		return
	}

	body := struct {
		Status string `json:"status"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "malformed request body, please provide valid JSON"})
		return
	}

	err = oh.UpdateOrderStatusByID(r.Context(), userID, orderID, body.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "no order with such an ID could be found"})
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not update the order status, something went wrong"})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "order status updated succesfully"})
}

func (oh *OrderHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.ClaimsFromContext(r.Context())
	userID := int(userClaims["sub"].(float64))

	order, err := oh.CheckoutOrder(r.Context(), userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not checkout the order, something went wrong"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, order)
}
