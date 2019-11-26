package main

import (
	"log"

	"github.com/aboglioli/big-brother/composition"
	infrComp "github.com/aboglioli/big-brother/infrastructure/composition"
	infrEvents "github.com/aboglioli/big-brother/infrastructure/events"
)

func main() {
	// Dendencies resolution
	eventMgr, err := infrEvents.GetManager()
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

	infrComp.StartREST(eventMgr, compositionService)
}
