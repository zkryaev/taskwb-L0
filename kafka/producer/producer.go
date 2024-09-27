package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/zkryaev/taskwb-L0/repository"
	"github.com/zkryaev/taskwb-L0/repository/config"
	"github.com/zkryaev/taskwb-L0/script"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	topic := "orders"

	cfg := config.Load(cfgPath)
	ordersRepo, err := repository.New(cfg)
	if err != nil {
		log.Fatal("Connection to DB failed", err)
	}
	defer ordersRepo.DB.Close()
	orders, err := ordersRepo.GetOrders()
	if err != nil {
		log.Fatal("Failed to get old orders from DB", err)
		return
	}
	log.Println("Producer is launched!")
	for {
		log.Println("Type 's' to generate order")
		log.Println("Type 'c' to send copy")
		log.Println("Type 'exit' to quit")
		var input string
		var orderJSON []byte
		fmt.Scanln(&input)

		if input == "exit" {
			fmt.Println("Exiting the program...")
			break
		}

		// send = сгенерировать
		if input == "s" {
			orderGenerated := script.GenerateOrder()
			orderJSON, err = json.Marshal(orderGenerated)
			if err != nil {
				log.Printf("Failed to convert order to JSON: %s", err)
				continue
			}
		}

		if input == "c" {
			log.Println("Choose 1 of old orders:")
			for i := 0; i < len(orders); i++ {
				fmt.Println(i, orders[i].OrderUID)
			}
			var indstr string
			fmt.Scanln(&indstr)
			ind, err := strconv.Atoi(indstr)
			if err != nil {
				log.Println("Entered is not a number!")
				continue
			}
			if ind < 0 || ind > len(orders) {
				log.Println("Entered number isn't in range of orders!")
				continue
			}
			orderJSON, err = json.Marshal(orders[ind])
			if err != nil {
				log.Printf("Failed to convert order to JSON: %s", err)
				continue
			}
		}

		err = PushOrderToQueue(topic, orderJSON)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %s", err)
			continue
		}

		log.Printf("Successfully sent order")
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
