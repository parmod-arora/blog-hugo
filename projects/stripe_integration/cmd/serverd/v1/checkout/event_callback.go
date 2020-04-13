package checkout

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"stripe-integration/internal/order"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

// EventCallback evant callback from stripe
func EventCallback(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// If you are testing your webhook locally with the Stripe CLI you
	// can find the endpoint's secret by running `stripe listen`
	// Otherwise, find your endpoint's secret in your webhook settings in the Developer Dashboard
	endpointSecret := os.Getenv("stripe_webhook_secret")

	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		order.UpdateSessionStatus(r.Context(), session.ID, order.StatusPaid)
	default:
		fmt.Fprintf(os.Stderr, "Unexpected event type: %s\\n", event.Type)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
