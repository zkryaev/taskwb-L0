package database

import (
	"database/sql"

	"github.com/zkryaev/taskwb-nats-stream/models"
)

func AddPayment(db *sql.DB, payment models.Payment, OrderUID string) error {
	query := `INSERT INTO "payments"("transaction", "request_id", "currency", "provider", "amount", "paymentdt", "bank", "delivery_cost", "goods_total", "custom_fee", "order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(
		query,
		payment.Transaction,
		payment.RequsetID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
		OrderUID,
	)
	if err != nil {
		return err
	}
	return nil
}
