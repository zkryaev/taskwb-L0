package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/zkryaev/taskwb-L0/script"
)

func main() {
	topic := "orders"
	log.Println("Producer is launched!")
	log.Println("Type 'exit' to quit:")
	for {
		// Спрашиваем у пользователя, что он хочет сделать
		var input string
		fmt.Scanln(&input)

		// Если пользователь ввел "exit", завершаем программу
		if input == "exit" {
			fmt.Println("Exiting the program...")
			break
		}

		// Генерация и отправка заказа
		order := script.GenerateOrder()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to convert order to JSON: %s", err)
			continue
		}

		err = PushOrderToQueue(topic, orderJSON)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %s", err)
			continue
		}

		log.Printf("Successfully sent order: %s", order.OrderUID)
	}
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	return sarama.NewSyncProducer(brokers, config)
}

func PushOrderToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:9092"}

	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Order is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topic,
		partition,
		offset,
	)

	return nil
}
