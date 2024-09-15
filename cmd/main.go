package main

import (
	"fmt"

	"github.com/zkryaev/taskwb-L0/cache"
	"github.com/zkryaev/taskwb-L0/database"
	"github.com/zkryaev/taskwb-L0/server"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	fmt.Println("Starting...")
	cfg := database.Load(cfgPath)
	db, err := database.Connect(cfg)
	defer db.Close()
	defer fmt.Println("DB: disconnected")
	if err != nil {
		fmt.Printf("Connection to DB is failed: %v\n", err)
		return
	}
	fmt.Println("DB: connected!")
	fmt.Printf("---POSTGRE---\nHost: %s\nPort: %s\nUser: %s\nName: %s\n", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name)

	cache := cache.New()
	orders, err := database.GetOrders(db)
	if err != nil {
		fmt.Println("failed to refil cache from db: %w", err)
		return
	}
	for _, order := range orders {
		cache.SaveOrder(order)
	}

	s := server.New(cfgPath, cache)
	s.Launch()
}

/*var objectNum uint = 5
orders := script.GenerateObjects(objectNum)
for i := range objectNum {
	err := database.AddOrder(db, orders[i])
	if err != nil {
		db.Close()
		log.Fatal(err)
		return
	}
}*/
