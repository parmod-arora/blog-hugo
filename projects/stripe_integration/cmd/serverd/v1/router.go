package v1

import (
	"stripe-integration/cmd/serverd/v1/checkout"

	"github.com/go-chi/chi"
)

// Router registers handlers to the router provided in the argument
func Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/create-session", checkout.CreateSession)
			r.Get("/session-status", checkout.OrderStatus)
		})
		r.Post("/webhooks", checkout.EventCallback)
	}
}
