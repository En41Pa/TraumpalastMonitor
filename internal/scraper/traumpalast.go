package scraper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Traumpalast struct {
	url       string
	userAgent string
	timeout   time.Duration
}

func NewTraumpalast(url, userAgent string, timeoutSeconds int) *Traumpalast {
	return &Traumpalast{
		url:       url,
		userAgent: userAgent,
		timeout:   time.Duration(timeoutSeconds) * time.Second,
	}
}

func (t *Traumpalast) FetchAvailableDates() ([]string, error) {
	client := &http.Client{
		Timeout: t.timeout,
	}

	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", t.userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return t.extractDates(string(body))
}

func (t *Traumpalast) extractDates(html string) ([]string, error) {
	var dates []string
	dateSet := make(map[string]bool)

	// Look for time elements with datetime attributes that indicate show times
	timePattern := regexp.MustCompile(`<time datetime="2025-(\d{2})-(\d{2})T\d{2}:\d{2}:\d{2}">(\d{2}:\d{2})</time>`)
	timeMatches := timePattern.FindAllStringSubmatch(html, -1)

	for _, match := range timeMatches {
		if len(match) >= 3 {
			month := match[1]
			day := match[2]

			// Convert to the format we want (DD.MM)
			dayInt, _ := strconv.Atoi(day)
			monthInt, _ := strconv.Atoi(month)
			dateStr := fmt.Sprintf("%d.%02d", dayInt, monthInt)

			if !dateSet[dateStr] {
				dates = append(dates, dateStr)
				dateSet[dateStr] = true
			}
		}
	}

	// Sort dates
	sort.Slice(dates, func(i, j int) bool {
		return t.compareDates(dates[i], dates[j])
	})

	return dates, nil
}

func (t *Traumpalast) compareDates(date1, date2 string) bool {
	d1 := t.parseDate(date1)
	d2 := t.parseDate(date2)
	return d1.Before(d2)
}

func (t *Traumpalast) parseDate(dateStr string) time.Time {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 2 {
		return time.Time{}
	}

	day, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])

	// Assume current year (2025 based on the HTML)
	return time.Date(2025, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
