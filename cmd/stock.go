package main

import (
	"log"

	"github.com/aboglioli/big-brother/infrastructure/rest"
)

func main() {
	log.Println("[INTERFACES] Starting...")

	rest.Start()

	log.Println("[INTERFACES] Started")
}
