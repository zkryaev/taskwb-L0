package main

import (
	"fmt"

	"github.com/zkryaev/taskwb-L0/database"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	fmt.Println("Starting...")
	cfg := database.Load(cfgPath)
	db, err := database.Connect(cfg)
	if err != nil {
		fmt.Printf("Failed main: %v\n", err)
		return
	}
	fmt.Println(db)
	fmt.Println("Successfull connected to DB")

}
