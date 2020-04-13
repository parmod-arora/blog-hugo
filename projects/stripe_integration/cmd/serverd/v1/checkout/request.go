package checkout

// CreateSessionRequest checkout request
type CreateSessionRequest struct {
	OrderID string       `json:"order_id"`
	Items   []OrderItems `json:"items"`
}

// OrderItems list of inventory items
type OrderItems struct {
	Amount      int64  `json:"amount"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
