package database

import (
	"database/sql"

	"github.com/zkryaev/taskwb-nats-stream/models"
)

func AddItems(db *sql.DB, items []models.Item, OrderUID string) (err error) {
	for _, item := range items {
		err = AddItem(db, item, OrderUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddItem(db *sql.DB, item models.Item, OrderUID string) error {
	query := `INSERT INTO "items"("chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status", "order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := db.Exec(
		query,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
		OrderUID,
	)
	if err != nil {
		return err
	}
	return nil
}
