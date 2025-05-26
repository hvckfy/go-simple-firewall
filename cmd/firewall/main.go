package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-simple-firewall/internal/config"
	"go-simple-firewall/internal/firewall"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем firewall
	fw, err := firewall.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create firewall: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down firewall...")
		fw.Shutdown()
		os.Exit(0)
	}()

	// Запускаем firewall
	log.Printf("🔥 Starting Go Simple Firewall...")
	if err := fw.Start(); err != nil {
		log.Fatalf("Failed to start firewall: %v", err)
	}
}
