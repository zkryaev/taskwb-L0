package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zkryaev/taskwb-L0/models"
)

func AddOrder(db *sql.DB, order models.Order) error {
	query := `INSERT INTO orders("order_uid", "track_number","entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard") 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(
		query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	err = AddPayment(db, order.Payment, order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	err = AddItems(db, order.Items, order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}

	err = AddDelivery(db, order.Delivery, order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}
	return nil
}

func GetOrders(db *sql.DB) ([]models.Order, error) {
	query := "SELECT * FROM orders"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		delivery, err := GetDelivery(db, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get delivery for order %s: %w", order.OrderUID, err)
		}
		order.Delivery = *delivery

		payment, err := GetPayment(db, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment for order %s: %w", order.OrderUID, err)
		}
		order.Payment = *payment

		items, err := GetItems(db, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get items for order %s: %w", order.OrderUID, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func GetOrder(db *sql.DB, OrderUID string) (*models.Order, error) {
	query := "SELECT * FROM orders WHERE order_uid = $1"
	row := db.QueryRow(query, OrderUID)
	var order models.Order
	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	delivery, err := GetDelivery(db, OrderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	payment, err := GetPayment(db, OrderUID)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := GetItems(db, OrderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}
