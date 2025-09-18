# Traumpalast Date Monitor 🎬

A Go application that monitors the Traumpalast Leonberg website for new movie screening dates and sends Discord notifications when new dates become available. Currently monitoring the "Demon Slayer: Kimetsu no Yaiba Infinity Castle" movie.

## Features

- **Web Scraping**: Uses Go's standard `net/http` library to scrape the Traumpalast website
- **Date Detection**: Automatically detects all available screening dates
- **New Date Monitoring**: Tracks when dates newer than 24.09.2025 become available
- **Discord Notifications**: Sends formatted notifications via Discord webhook
- **Test Mode**: Allows testing for specific dates
- **Continuous Monitoring**: Runs continuously with configurable intervals

## Setup

### 1. Configuration

The application uses a `config.json` file for all settings. The file is already included with your Discord webhook configured:

```json
{
  "discord_webhook_url": "https://discord.com/api/webhooks/1418140603568623707/VzVkNdAxpfFK_7VxAmWroSGAttwDKxoFmYElRC0GYs5hq6w0eDKh4SGX2skL5f9AX6Wr",
  "traumpalast_url": "https://leonberg.traumpalast.de/index.php/FN/8182/PID/5574.html",
  "current_newest_date": "24.09",
  "check_interval_minutes": 30,
  "request_timeout_seconds": 30,
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
}
```

You can modify any of these settings:
- **discord_webhook_url**: Your Discord webhook URL
- **traumpalast_url**: The movie page to monitor
- **current_newest_date**: Reference date to compare against
- **check_interval_minutes**: How often to check (default: 30 minutes)
- **request_timeout_seconds**: HTTP request timeout
- **user_agent**: Browser user agent string

### 2. Build and Run

```bash
# Build the application
go build -o traumpalast-monitor main.go

# Or run directly
go run main.go
```

## Usage

### Continuous Monitoring Mode (Default)

```bash
./traumpalast-monitor
# or
go run main.go
```

This will:
- Check for new dates every 30 minutes
- Send Discord notifications when dates after 24.09.2025 are found
- Run indefinitely until stopped

### Single Check Mode

```bash
./traumpalast-monitor check
# or
go run main.go check
```

Performs a single check and exits.

### Test Mode

Test if a specific date is currently available:

```bash
./traumpalast-monitor test "25.09"
# or
go run main.go test "25.09"
```

This will:
- Search for the specified date on the website
- Send a Discord notification with the test result
- Show whether the date was found or not

## Example Outputs

### When New Dates Are Found

```
🎬 **New dates available at Traumpalast IMAX!** 🎬

New latest date: **25.09.2025**
Previous latest date was: **24.09.2025**

Check it out: https://leonberg.traumpalast.de/index.php/FN/8182/PID/5574.html
```

### Test Results

```
🎯 **Test Result: SUCCESS!** 🎯

The date **24.09.2025** was found on the Traumpalast website!

URL: https://leonberg.traumpalast.de/index.php/FN/8182/PID/5574.html
```

## Configuration

All settings are managed through the `config.json` file:

- **discord_webhook_url**: Discord webhook for notifications
- **traumpalast_url**: Movie page URL to monitor  
- **current_newest_date**: Reference date (currently "24.09")
- **check_interval_minutes**: Monitoring frequency (30 minutes)
- **request_timeout_seconds**: HTTP timeout (30 seconds)
- **user_agent**: Browser identification string

## How It Works

1. **Web Scraping**: Fetches the HTML from the Traumpalast website
2. **Date Extraction**: Uses regex patterns to find dates in format "DD.MM."
3. **Date Parsing**: Converts dates to Go `time.Time` for comparison
4. **Monitoring**: Compares found dates against the reference date (24.09)
5. **Notifications**: Sends Discord webhook when new dates are detected

## Error Handling

- Network errors are logged and monitoring continues
- Discord webhook failures are logged but don't stop monitoring
- Invalid dates are ignored
- HTTP errors are reported with status codes

## Architecture

The application follows a clean, modular architecture:

```
internal/
├── config/     # Configuration management
├── monitor/    # Date monitoring orchestration  
├── notifier/   # Discord webhook notifications
└── scraper/    # Web scraping functionality
```

### Package Responsibilities:
- **config**: Loads and validates JSON configuration
- **scraper**: Fetches and parses Traumpalast website data
- **notifier**: Handles Discord webhook notifications
- **monitor**: Orchestrates scraping, date comparison, and notifications
- **main**: CLI interface and application lifecycle

## Dependencies

Uses only Go standard library:
- `net/http` for web requests
- `regexp` for date extraction
- `time` for date parsing and scheduling
- `encoding/json` for configuration and webhook payloads

## Monitoring Tips

- Run as a background service or in a screen/tmux session
- Monitor logs for any issues
- Test the Discord webhook before starting continuous monitoring
- The application will continue running even if individual checks fail
