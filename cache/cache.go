package cache

import "github.com/zkryaev/taskwb-L0/models"

type Cache struct {
	orders map[string]models.Order
}

func New() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

func (c *Cache) SaveOrder(order models.Order) {
	c.orders[order.OrderUID] = order
}

func (c *Cache) GetOrder(OrderUID string) (models.Order, bool) {
	order, ok := c.orders[OrderUID]
	return order, ok
}
