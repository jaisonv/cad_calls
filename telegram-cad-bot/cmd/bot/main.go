package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaisonv/telegram-cad-bot/internal/bot"
	"github.com/jaisonv/telegram-cad-bot/internal/scheduler"
)

func main() {
	// Command line flags
	token := flag.String("token", os.Getenv("TELEGRAM_BOT_TOKEN"), "Telegram bot token")
	baseURL := flag.String("base-url", "https://southmiamipdfl.policetocitizen.com", "CAD API base URL")
	agencyID := flag.Int("agency-id", 386, "Agency ID")
	dbPath := flag.String("db", "./cadbot.db", "Database file path")
	checkInterval := flag.Int("interval", 5, "Check interval in minutes")
	verifySSL := flag.Bool("verify-ssl", false, "Verify SSL certificates")
	flag.Parse()

	if *token == "" {
		log.Fatal("Telegram bot token is required. Set TELEGRAM_BOT_TOKEN environment variable or use -token flag")
	}

	logger := log.New(os.Stdout, "[CADBot] ", log.LstdFlags)

	// Create bot configuration
	config := &bot.Config{
		TelegramToken: *token,
		BaseURL:       *baseURL,
		AgencyID:      *agencyID,
		VerifySSL:     *verifySSL,
		Timeout:       30 * time.Second,
		DBPath:        *dbPath,
	}

	// Initialize bot
	logger.Println("Initializing bot...")
	b, err := bot.NewBot(config, logger)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Initialize poller
	logger.Println("Initializing poller...")
	poller := scheduler.NewPoller(b, b.GetDB(), time.Duration(*checkInterval)*time.Minute, logger)

	// Start poller
	poller.Start()

	// Start bot in a goroutine
	go b.Start()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Println("Bot is running. Press Ctrl+C to stop.")
	<-sigChan

	// Graceful shutdown
	logger.Println("Shutting down...")
	poller.Stop()
	b.Stop()
	logger.Println("Goodbye!")
}
