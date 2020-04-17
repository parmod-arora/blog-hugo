package checkout

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stripe-integration/internal/order"
	"stripe-integration/pkg/web"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

// CreateSession Route for creating session for stripe checkout
func CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := CreateSessionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	model := transformToModel(req)

	model, err := order.CreatePendingOrder(ctx, model)
	if err != nil {
		log.Printf("error CreatePendingOrder %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set your secret key. Remember to switch to your live secret key in production!
	// See your keys here: https://dashboard.stripe.com/account/apikeys
	stripe.Key = os.Getenv("stripe_secret")
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SuccessURL: stripe.String("http://localhost:3000/client/success.html?session_id={CHECKOUT_SESSION_ID}&order_id=" + model.OrderID),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
	}
	for _, orderItem := range model.Items {
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(orderItem.Name),
			Description: stripe.String(orderItem.Description),
			Amount:      stripe.Int64(orderItem.Amount),
			Currency:    stripe.String(string(stripe.CurrencySGD)),
			Quantity:    stripe.Int64(1),
		})
	}
	session, err := session.New(params)
	if err != nil {
		log.Printf("error create session %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := order.UpdateOrderStatus(ctx, model.ID, session.ID, order.StatusCreated); err != nil {
		log.Printf("error update Order %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.Respond(w, http.StatusOK, session, nil)
}

func transformToModel(req CreateSessionRequest) order.Order {
	model := order.Order{
		OrderID: req.OrderID,
	}
	for _, item := range req.Items {
		model.Items = append(model.Items, order.Item{
			Amount:      item.Amount,
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return model
}

type orderStatusResponse struct {
	OrderID   string `json:"order_id"`
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
}

// OrderStatus retrieve order status
func OrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.URL.Query().Get("order_id")
	order, err := order.RetrieveOrderByID(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := orderStatusResponse{
		OrderID:   order.OrderID,
		SessionID: *order.SessionID,
		Status:    string(order.Status),
	}
	web.Respond(w, http.StatusOK, response, nil)
}
