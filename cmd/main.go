package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zkryaev/taskwb-L0/database"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	fmt.Println("Starting...")
	cfg := database.Load(cfgPath)
	fmt.Printf("---PGINFO---\nHost: %s\nPort: %s\nUser: %s\nName: %s\n------------\n", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name)
	db, err := database.Connect(cfg)
	if err != nil {
		fmt.Printf("Failed main: %v\n", err)
		return
	}
	fmt.Println("Connected to DB!")
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

	order, err := database.GetOrder(db, "kIMyUMvybsmE")
	if err != nil {
		db.Close()
		log.Fatal(err)
		return
	}
	orderJSON, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal order to JSON: %v", err)
	}

	// Вывод JSON на экран
	fmt.Println(string(orderJSON))
	db.Close()
	fmt.Println("Close connection to DB")
}
