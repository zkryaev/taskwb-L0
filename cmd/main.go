package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/zkryaev/taskwb-L0/cache"
	"github.com/zkryaev/taskwb-L0/controller"
	"github.com/zkryaev/taskwb-L0/kafka/consumer"
	"github.com/zkryaev/taskwb-L0/repository"
	"github.com/zkryaev/taskwb-L0/repository/config"
	"go.uber.org/zap"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application...")

	// Загружаем конфигурацию
	cfg := config.Load(cfgPath)

	// Подключаемся к базе данных
	ordersRepo, err := repository.New(cfg)
	if err != nil {
		logger.Error("Connection to DB failed", zap.Error(err))
		return
	}
	defer ordersRepo.DB.Close()
	logger.Info(
		"DB connected successfully",
		zap.String("host", cfg.DB.Host),
		zap.String("port", cfg.DB.Port),
		zap.String("db", cfg.DB.Name),
		zap.String("user", cfg.DB.User),
	)

	// Инициализируем кэш
	cache := cache.New()

	// Получаем заказы из базы данных и заполняем кэш
	orders, err := ordersRepo.GetOrders()
	if err != nil {
		logger.Error("Failed to refill cache from DB", zap.Error(err))
		return
	}

	logger.Info("Restoring the cache...")
	for _, order := range orders {
		cache.SaveOrder(order)
		logger.Info("Cached", zap.String("order_uid", order.OrderUID))
	}

	// Запуск сервера
	s := controller.New(cfgPath, cache)
	go func() {
		s.Launch()
	}()
	logger.Info(
		"Server launched",
		zap.String("host", cfg.App.Host),
		zap.String("port", cfg.App.Port),
	)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	//connect consumer
	err = consumer.Subscribe(cache, ordersRepo, *logger, sigchan)
	if err != nil {
		logger.Fatal("consumer error:", zap.Error(err))
	}
	logger.Info("Application shutdown")
}
