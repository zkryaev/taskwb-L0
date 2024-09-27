package consumer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/zkryaev/taskwb-L0/cache"
	"github.com/zkryaev/taskwb-L0/models"
	"github.com/zkryaev/taskwb-L0/repository"
	"go.uber.org/zap"
)

func Subscribe(cache *cache.Cache, db *repository.OrdersRepo, logger zap.Logger, sigchan <-chan os.Signal) error {
	topic := "orders"

	worker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		return fmt.Errorf("connect_consumer failed: %w", err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("consume_partition failed: %w", err)
	}

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				logger.Error("consuming error:", zap.Error(err))
			case msg := <-consumer.Messages():
				var order models.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					logger.Error("failed to unmarshal message", zap.Error(err))
					continue
				}
				if _, found := cache.GetOrder(order.OrderUID); found {
					logger.Info("order exist, didn't add")
					continue
				}
				if err := db.AddOrder(order); err != nil {
					logger.Error("failed to save order to DB", zap.Error(err))
					continue
				}
				cache.SaveOrder(order)
				logger.Info("Consumed", zap.String("order_uid", order.OrderUID))
			case <-sigchan:
				doneCh <- struct{}{}
			}
		}
	}()
	logger.Info("Consumer subscribed to Kafka!")
	<-doneCh
	close(doneCh)
	if err := worker.Close(); err != nil {
		return err
	}
	return nil
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}
