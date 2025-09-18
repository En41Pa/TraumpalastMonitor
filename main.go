package main

import (
	"log"
	"os"
	"time"

	"github.com/en41pa/TraumpalastMonitor/internal/config"
	"github.com/en41pa/TraumpalastMonitor/internal/monitor"
)

func main() {
	log.Println("🎬 Traumpalast Date Monitor Starting...")

	// Load configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration successfully")
	log.Printf("Monitoring URL: %s", cfg.TraumpalastURL)
	log.Printf("Check interval: %d minutes", cfg.CheckIntervalMinutes)

	if cfg.DiscordWebhookURL == "" {
		log.Println("⚠️  No Discord webhook URL configured. Notifications will be logged only.")
	} else {
		log.Println("✅ Discord webhook configured")
	}

	dateMonitor := monitor.New(cfg)

	// Check command line arguments for different modes
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			if len(os.Args) > 2 {
				// Test mode: check for a specific date
				testDate := os.Args[2]
				dateMonitor.TestSpecificDate(testDate)
			} else {
				log.Println("Usage: go run main.go test <date>")
				log.Println("Example: go run main.go test \"24.09\"")
			}
			return
		case "check":
			// Single check mode
			err := dateMonitor.CheckForNewDates()
			if err != nil {
				log.Printf("Error during check: %v", err)
			}
			return
		default:
			log.Printf("Unknown command: %s", os.Args[1])
			log.Println("Available commands:")
			log.Println("  test <date>  - Test if a specific date is available")
			log.Println("  check        - Run a single check and exit")
			log.Println("  (no args)    - Run continuous monitoring")
			return
		}
	}

	// Continuous monitoring mode
	log.Println("Starting continuous monitoring mode...")
	log.Printf("Monitoring for dates after: %s", cfg.CurrentNewestDate)

	// Send startup notification
	if cfg.DiscordWebhookURL != "" {
		err = dateMonitor.SendStartupNotification()
		if err != nil {
			log.Printf("Failed to send startup notification: %v", err)
		} else {
			log.Println("Startup notification sent to Discord!")
		}
	}

	// Initial check
	err = dateMonitor.CheckForNewDates()
	if err != nil {
		log.Printf("Initial check failed: %v", err)
	}

	// Set up ticker for regular checks
	ticker := time.NewTicker(time.Duration(cfg.CheckIntervalMinutes) * time.Minute)
	defer ticker.Stop()

	log.Printf("Monitoring every %d minutes. Press Ctrl+C to stop.", cfg.CheckIntervalMinutes)

	for {
		select {
		case <-ticker.C:
			err := dateMonitor.CheckForNewDates()
			if err != nil {
				log.Printf("Monitoring check failed: %v", err)
			}
		}
	}
}
