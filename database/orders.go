package database

import (
	"database/sql"
	"fmt"

	"github.com/zkryaev/taskwb-nats-stream/models"
)

func AddOrder(db *sql.DB, order models.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	query := `INSERT INTO orders("order_uid", "track_number","entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard") 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(
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
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	err = AddPayment(db, order.Payment, order.OrderUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	err = AddItems(db, order.Items, order.OrderUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert items: %w", err)
	}

	err = AddDelivery(db, order.Delivery, order.OrderUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert delivery: %w", err)
	}
	return nil
}
