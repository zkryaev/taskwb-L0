package database

import (
	"database/sql"

	"github.com/zkryaev/taskwb-nats-stream/models"
)

func AddDelivery(db *sql.DB, delivery models.Delivery, OrderUID string) error {
	query := `INSERT INTO deliveries ("name", "phone", "zip", "city", "address", "region", "email", "order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(
		query,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		OrderUID,
	)
	if err != nil {
		return err
	}
	return nil
}
