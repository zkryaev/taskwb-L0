package repository

import (
	"github.com/zkryaev/taskwb-L0/models"
)

type Orders interface {
	AddOrder(order models.Order) error
	GetOrder(OrderUID string) (*models.Order, error)
	GetOrders() ([]models.Order, error)
}
