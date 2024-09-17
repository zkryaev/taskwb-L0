package cache

import (
	"sync"

	"github.com/zkryaev/taskwb-L0/models"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func New() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

func (c *Cache) SaveOrder(order models.Order) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) GetOrder(OrderUID string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	order, ok := c.orders[OrderUID]
	return order, ok
}
