package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Discord struct {
	webhookURL string
}

type webhookPayload struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Thumbnail struct {
	URL string `json:"url,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

const (
	ColorGreen  = 0x00ff00 // Success/Found
	ColorRed    = 0xff0000 // Error/Not Found
	ColorBlue   = 0x0099ff // Info/New Dates
	ColorOrange = 0xff9900 // Warning
)

func NewDiscord(webhookURL string) *Discord {
	return &Discord{
		webhookURL: webhookURL,
	}
}

func (d *Discord) SendNotification(payload webhookPayload) error {
	if d.webhookURL == "" {
		log.Println("No webhook URL configured, skipping Discord notification")
		return nil
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook data: %v", err)
	}

	resp, err := http.Post(d.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (d *Discord) SendNewDateNotification(newDate, previousDate, url string) error {
	embed := Embed{
		Title:       "🎬 New Movie Dates Available!",
		Description: "**Demon Slayer: Infinity Castle** has new screening dates at Traumpalast Leonberg!",
		Color:       ColorBlue,
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: Thumbnail{
			URL: "https://images.unsplash.com/photo-1489599735734-79b4212eea45?w=400&h=400&fit=crop&crop=center",
		},
		Fields: []Field{
			{
				Name:   "🆕 New Latest Date",
				Value:  fmt.Sprintf("**%s.2025**", newDate),
				Inline: true,
			},
			{
				Name:   "📅 Previous Latest Date",
				Value:  fmt.Sprintf("**%s.2025**", previousDate),
				Inline: true,
			},
			{
				Name:   "🎯 Action Required",
				Value:  "Book your tickets now before they sell out!",
				Inline: false,
			},
			{
				Name:   "🎪 Cinema",
				Value:  "**Traumpalast Leonberg IMAX**",
				Inline: true,
			},
			{
				Name:   "🎭 Movie",
				Value:  "**Demon Slayer: Infinity Castle**",
				Inline: true,
			},
		},
		Footer: Footer{
			Text:    "Traumpalast Monitor • Powered by Go",
			IconURL: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=32&h=32&fit=crop&crop=center",
		},
	}

	payload := webhookPayload{
		Content: fmt.Sprintf("🚨 **ALERT:** New dates detected! 🚨\n[**📖 Book Tickets Here**](%s)", url),
		Embeds:  []Embed{embed},
	}

	return d.SendNotification(payload)
}

func (d *Discord) SendTestResult(targetDate string, found bool, availableDates []string, url string) error {
	var embed Embed
	var content string

	if found {
		embed = Embed{
			Title:       "✅ Test Result: SUCCESS!",
			Description: fmt.Sprintf("The date **%s.2025** was found on the Traumpalast website!", targetDate),
			Color:       ColorGreen,
			Timestamp:   time.Now().Format(time.RFC3339),
			Thumbnail: Thumbnail{
				URL: "https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=400&h=400&fit=crop&crop=center",
			},
			Fields: []Field{
				{
					Name:   "🎯 Target Date",
					Value:  fmt.Sprintf("**%s.2025**", targetDate),
					Inline: true,
				},
				{
					Name:   "📊 Status",
					Value:  "**Available** ✅",
					Inline: true,
				},
				{
					Name:   "🎬 Next Step",
					Value:  "This date has screenings available!",
					Inline: false,
				},
			},
			Footer: Footer{
				Text:    "Test completed successfully",
				IconURL: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=32&h=32&fit=crop&crop=center",
			},
		}
		content = fmt.Sprintf("🎉 **Test completed!** [**View Screenings**](%s)", url)
	} else {
		// Format available dates nicely
		datesStr := ""
		for i, date := range availableDates {
			if i > 0 {
				datesStr += ", "
			}
			datesStr += fmt.Sprintf("**%s.2025**", date)
		}
		if datesStr == "" {
			datesStr = "No dates found"
		}

		embed = Embed{
			Title:       "❌ Test Result: NOT FOUND",
			Description: fmt.Sprintf("The date **%s.2025** was not found on the Traumpalast website.", targetDate),
			Color:       ColorRed,
			Timestamp:   time.Now().Format(time.RFC3339),
			Thumbnail: Thumbnail{
				URL: "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=400&h=400&fit=crop&crop=center",
			},
			Fields: []Field{
				{
					Name:   "🎯 Target Date",
					Value:  fmt.Sprintf("**%s.2025**", targetDate),
					Inline: true,
				},
				{
					Name:   "📊 Status",
					Value:  "**Not Available** ❌",
					Inline: true,
				},
				{
					Name:   "📅 Available Dates",
					Value:  datesStr,
					Inline: false,
				},
			},
			Footer: Footer{
				Text:    "Keep monitoring for updates",
				IconURL: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=32&h=32&fit=crop&crop=center",
			},
		}
		content = fmt.Sprintf("🔍 **Test completed!** [**Check Website**](%s)", url)
	}

	payload := webhookPayload{
		Content: content,
		Embeds:  []Embed{embed},
	}

	return d.SendNotification(payload)
}

func (d *Discord) SendStartupNotification(checkInterval int, referenceDate, url string) error {
	embed := Embed{
		Title:       "🚀 Traumpalast Monitor Started",
		Description: "Your movie date monitor is now actively watching for new screening dates!",
		Color:       ColorGreen,
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: Thumbnail{
			URL: "https://images.unsplash.com/photo-1489599735734-79b4212eea45?w=400&h=400&fit=crop&crop=center",
		},
		Fields: []Field{
			{
				Name:   "🎬 Monitoring",
				Value:  "**Demon Slayer: Infinity Castle**",
				Inline: true,
			},
			{
				Name:   "🏛️ Cinema",
				Value:  "**Traumpalast Leonberg**",
				Inline: true,
			},
			{
				Name:   "⏰ Check Interval",
				Value:  fmt.Sprintf("Every **%d minutes**", checkInterval),
				Inline: true,
			},
			{
				Name:   "📅 Reference Date",
				Value:  fmt.Sprintf("Watching for dates after **%s.2025**", referenceDate),
				Inline: true,
			},
			{
				Name:   "🎯 Status",
				Value:  "**Active & Monitoring** 🟢",
				Inline: false,
			},
		},
		Footer: Footer{
			Text:    "You'll be notified when new dates become available",
			IconURL: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=32&h=32&fit=crop&crop=center",
		},
	}

	payload := webhookPayload{
		Content: fmt.Sprintf("🎭 **Monitor activated!** I'll watch for new dates and notify you immediately.\n[**Current Screenings**](%s)", url),
		Embeds:  []Embed{embed},
	}

	return d.SendNotification(payload)
}
