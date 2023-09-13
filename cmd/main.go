package main

import (
	"log"

	"lab1/internal/api"
)

func main() {
	log.Println("Application Start")
	api.StartServer()
	log.Println("Application terminated!")
}
