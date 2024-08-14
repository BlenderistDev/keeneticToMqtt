package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"keeneticToMqtt/internal/app"
)

func main() {
	cont, err := app.NewContainer()
	if err != nil {
		panic(fmt.Errorf("error while creating container: %w", err))
	}

	done := cont.EntityManager.Run()

	sig := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	shutdownCh := make(chan os.Signal, len(sig))
	signal.Notify(shutdownCh, sig...)

	<-shutdownCh
	done <- true
	cont.Logger.Info("process interrupted by signal")
	return
}
