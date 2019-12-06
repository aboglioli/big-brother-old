package main

import (
	"log"

	"github.com/aboglioli/big-brother/composition"
	implComp "github.com/aboglioli/big-brother/impl/composition"
	"github.com/aboglioli/big-brother/impl/events"
)

func main() {
	// Dendencies resolution
	eventMgr, err := events.Rabbit()
	if err != nil {
		log.Fatal(err)
		return
	}

	compositionRepository, err := composition.NewRepository()
	if err != nil {
		log.Fatal(err)
		return
	}

	compositionService := composition.NewService(compositionRepository, eventMgr)

	implComp.StartREST(eventMgr, compositionService)
}
