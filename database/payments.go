package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zkryaev/taskwb-L0/models"
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

func GetPayment(db *sql.DB, OrderUID string) (*models.Payment, error) {
	query := "SELECT * FROM payments WHERE order_uid = $1"
	row := db.QueryRow(query, OrderUID)
	var payment models.Payment
	err := row.Scan(&payment.Transaction, &payment.RequsetID, &payment.Provider, &payment.PaymentDT, &payment.GoodsTotal, &payment.DeliveryCost, &payment.CustomFee, &payment.Currency, &payment.Bank, &payment.Amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("get payment failed: %w", err)
	}
	return &payment, nil
}
