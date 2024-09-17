package script

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zkryaev/taskwb-L0/models"
)

/* Этот код в main

var objectNum uint = 1
	MockOrders := script.GenerateObjects(objectNum)
	for i := range objectNum {
		err = ordersRepo.AddOrder(MockOrders[i])
		if err != nil {
			ordersRepo.DB.Close()
			log.Fatal(err)
			return
		}
	}

*/

func GenerateObjects(cnt uint) []models.Order {
	rand.Seed(time.Now().UnixNano()) // Устанавливаем начальное значение для генератора случайных чисел
	var orders []models.Order
	for range cnt {
		delivery := models.Delivery{
			Name:    randomString(10),
			Phone:   randomPhone(),
			Zip:     randomZip(),
			City:    randomString(8),
			Address: randomString(15),
			Region:  randomString(8),
			Email:   randomString(5),
		}

		item := models.Item{
			ChrtID:      rand.Intn(1000),
			TrackNumber: randomString(10),
			Price:       rand.Intn(1000),
			Rid:         randomString(6),
			Name:        randomString(10),
			Sale:        rand.Intn(100),
			Size:        randomSize(),
			TotalPrice:  rand.Intn(1000),
			NmID:        rand.Intn(1000),
			Brand:       randomString(8),
			Status:      rand.Intn(5),
		}

		payment := models.Payment{
			Transaction:  randomString(10),
			RequestID:    randomString(8),
			Currency:     "USD",
			Provider:     randomString(6),
			Amount:       rand.Intn(10000),
			PaymentDT:    int(time.Now().Unix()),
			Bank:         randomString(6),
			DeliveryCost: rand.Intn(500),
			GoodsTotal:   rand.Intn(10000),
			CustomFee:    rand.Intn(100),
		}

		order := models.Order{
			OrderUID:          randomString(12),
			TrackNumber:       randomString(10),
			Entry:             randomString(5),
			Delivery:          delivery,
			Payment:           payment,
			Items:             []models.Item{item},
			Locale:            "en",
			InternalSignature: randomString(8),
			CustomerID:        randomString(8),
			DeliveryService:   randomString(5),
			Shardkey:          randomString(5),
			SmID:              rand.Intn(100),
			DateCreated:       time.Now().Format("2006-01-02"),
			OofShard:          randomString(4),
		}
		orders = append(orders, order)
	}
	return orders
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Функция для генерации случайного номера телефона
func randomPhone() string {
	return fmt.Sprintf("+1%010d", rand.Int63n(10000000000))
}

// Функция для генерации случайного ZIP-кода
func randomZip() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

// Функция для генерации случайного размера
func randomSize() string {
	sizes := []string{"S", "M", "L", "XL", "XXL"}
	return sizes[rand.Intn(len(sizes))]
}
