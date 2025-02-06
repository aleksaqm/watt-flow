package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	amqpURI := os.Getenv("AMQP_URI")
	emailSecret := os.Getenv("EMAIL_SECRET")
	influxURI := os.Getenv("INFLUX_URI")
	influxToken := os.Getenv("INFLUX_TOKEN")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_PORT := os.Getenv("DB_PORT")
	DB_TABLE := os.Getenv("DB_TABLE")

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_TABLE,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cleanup on exit

	worker, err := NewWorker(amqpURI, influxURI, influxToken, pgConnStr, emailSecret)
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}

	if err := worker.ConnectToBroker(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}

	if err := worker.Start(ctx); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal")

	// Perform graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer shutdownCancel()

	if err := worker.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Worker shut down successfully.")
}

func calculatePrice(spentPower float64, pricelist Pricelist) float64 {
	const billingPowerConstant = 7.0

	greenConsumption := min(spentPower, 350)
	blueConsumption := min(max(spentPower-350, 0), 1250)
	redConsumption := max(spentPower-1600, 0)

	basePrice := (greenConsumption * pricelist.GreenZone) +
		(blueConsumption * pricelist.BlueZone) +
		(redConsumption * pricelist.RedZone) +
		(billingPowerConstant * pricelist.BillingPower)

	// Apply tax
	finalPrice := basePrice + (basePrice * pricelist.Tax / 100)

	return finalPrice
}
