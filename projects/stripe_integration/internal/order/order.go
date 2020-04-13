package order

import (
	"context"
	"stripe-integration/internal/db"
	"time"
)

// Status order status
type Status string

// StatusPending pending order status
const StatusPending Status = "P"

// StatusPaid paid order status
const StatusPaid Status = "S"

// StatusCreated created order status
const StatusCreated Status = "C"

// StatusFailure failure order status
const StatusFailure Status = "F"

// Order structure
type Order struct {
	ID        int64
	OrderID   string
	SessionID *string
	Status    Status
	Items     []Item
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// Item order item
type Item struct {
	ID          int64
	OrderID     string
	Name        string
	Description string
	Amount      int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// CreatePendingOrder save order in DB with pending status
func CreatePendingOrder(ctx context.Context, order Order) (Order, error) {
	dbconn := db.GetConnection()
	const q = "INSERT INTO checkout_order (order_id, status) VALUES ($1, $2) RETURNING id, created_at, updated_at"
	tx, err := dbconn.Begin()
	if err != nil {
		return order, err
	}
	defer tx.Rollback()
	if err := dbconn.QueryRowContext(ctx, q, order.OrderID, StatusPending).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return order, err
	}

	for _, item := range order.Items {
		const q = "INSERT INTO checkout_order_item (order_id, name, description, amount) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at"
		if err := dbconn.QueryRowContext(ctx, q, order.ID, item.Name, item.Description, item.Amount).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return order, err
		}
	}
	return order, tx.Commit()
}

// RetrieveOrderByID retrieve order By order_id
func RetrieveOrderByID(ctx context.Context, orderID string) (Order, error) {
	var order Order
	dbconn := db.GetConnection()
	q := "SELECT id, order_id, session_id, status, created_at, updated_at from checkout_order WHERE order_id=$1"
	if err := dbconn.QueryRowContext(ctx, q, orderID).Scan(&order.ID, &order.OrderID, &order.SessionID, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return order, err
	}

	q = "SELECT id, order_id, name, description, amount, created_at, updated_at from checkout_order_item WHERE order_id=$1"
	rows, err := dbconn.QueryContext(ctx, q, order.ID)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.OrderID, &item.Name, &item.Description, &item.Amount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}

// UpdateOrderStatus update order status in DB
func UpdateOrderStatus(ctx context.Context, orderID int64, sessionID string, status Status) error {
	dbconn := db.GetConnection()
	const q = "UPDATE checkout_order set status=$1, session_id=$3 WHERE id=$2"
	_, err := dbconn.ExecContext(ctx, q, status, orderID, sessionID)
	return err
}

// UpdateSessionStatus update order status in DB
func UpdateSessionStatus(ctx context.Context, sessionID string, status Status) error {
	dbconn := db.GetConnection()
	const q = "UPDATE checkout_order set status=$1 WHERE session_id=$2"
	_, err := dbconn.ExecContext(ctx, q, status, sessionID)
	return err
}
