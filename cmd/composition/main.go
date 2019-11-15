package main

import (
	"github.com/aboglioli/big-brother/infrastructure/composition"
	"github.com/aboglioli/big-brother/infrastructure/events"
)

func main() {
	events.StartListeners()
	composition.StartREST()
}
