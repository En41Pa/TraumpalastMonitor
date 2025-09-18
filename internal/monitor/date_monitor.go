package monitor

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/en41pa/TraumpalastMonitor/internal/config"
	"github.com/en41pa/TraumpalastMonitor/internal/notifier"
	"github.com/en41pa/TraumpalastMonitor/internal/scraper"
)

type DateMonitor struct {
	config        *config.Config
	scraper       *scraper.Traumpalast
	notifier      *notifier.Discord
	lastKnownDate string
}

func New(cfg *config.Config) *DateMonitor {
	return &DateMonitor{
		config:        cfg,
		scraper:       scraper.NewTraumpalast(cfg.TraumpalastURL, cfg.UserAgent, cfg.RequestTimeoutSeconds),
		notifier:      notifier.NewDiscord(cfg.DiscordWebhookURL),
		lastKnownDate: cfg.CurrentNewestDate,
	}
}

func (dm *DateMonitor) SendStartupNotification() error {
	return dm.notifier.SendStartupNotification(
		dm.config.CheckIntervalMinutes,
		dm.config.CurrentNewestDate,
		dm.config.TraumpalastURL,
	)
}

func (dm *DateMonitor) CheckForNewDates() error {
	log.Println("Checking for new dates...")

	dates, err := dm.scraper.FetchAvailableDates()
	if err != nil {
		return fmt.Errorf("failed to scrape dates: %v", err)
	}

	log.Printf("Found dates: %v", dates)

	// Find the newest date
	var newestDate string
	for _, date := range dates {
		if newestDate == "" || dm.isDateAfter(date, newestDate) {
			newestDate = date
		}
	}

	if newestDate == "" {
		log.Println("No dates found")
		return nil
	}

	log.Printf("Current newest date: %s, Last known date: %s", newestDate, dm.lastKnownDate)

	// Check if there's a new date after our reference date
	if dm.isDateAfter(newestDate, dm.lastKnownDate) {
		err := dm.notifier.SendNewDateNotification(newestDate, dm.lastKnownDate, dm.config.TraumpalastURL)
		if err != nil {
			log.Printf("Failed to send Discord notification: %v", err)
		} else {
			log.Println("Discord notification sent successfully!")
		}

		// Update our reference
		dm.lastKnownDate = newestDate
	} else {
		log.Printf("No new dates found. Latest date is still %s", newestDate)
	}

	return nil
}

func (dm *DateMonitor) TestSpecificDate(targetDate string) {
	log.Printf("Testing for specific date: %s", targetDate)

	dates, err := dm.scraper.FetchAvailableDates()
	if err != nil {
		log.Printf("Error scraping dates: %v", err)
		return
	}

	log.Printf("All found dates: %v", dates)

	found := false
	for _, date := range dates {
		if date == targetDate {
			found = true
			break
		}
	}

	if found {
		log.Printf("✅ Date %s was found!", targetDate)
	} else {
		log.Printf("❌ Date %s was NOT found.", targetDate)
	}

	err = dm.notifier.SendTestResult(targetDate, found, dates, dm.config.TraumpalastURL)
	if err != nil {
		log.Printf("Failed to send test notification: %v", err)
	}
}

func (dm *DateMonitor) isDateAfter(dateStr, referenceDate string) bool {
	date := dm.parseDate(dateStr)
	reference := dm.parseDate(referenceDate)
	return date.After(reference)
}

func (dm *DateMonitor) parseDate(dateStr string) time.Time {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 2 {
		return time.Time{}
	}

	day, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])

	// Assume current year (2025 based on the HTML)
	return time.Date(2025, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
