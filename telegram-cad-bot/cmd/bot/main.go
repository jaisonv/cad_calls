package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jaisonv/telegram-cad-bot/internal/bot"
	"github.com/jaisonv/telegram-cad-bot/internal/scheduler"
)

func main() {
	// Command line flags
	token := flag.String("token", os.Getenv("TELEGRAM_BOT_TOKEN"), "Telegram bot token")
	pythonScript := flag.String("python-script", "", "Path to direct_api_post.py script")
	dbPath := flag.String("db", "./cadbot.db", "Database file path")
	checkInterval := flag.Int("interval", 5, "Check interval in minutes")
	flag.Parse()

	if *token == "" {
		log.Fatal("Telegram bot token is required. Set TELEGRAM_BOT_TOKEN environment variable or use -token flag")
	}

	logger := log.New(os.Stdout, "[CADBot] ", log.LstdFlags)

	// Determine Python script path
	scriptPath := *pythonScript
	if scriptPath == "" {
		// Default: assume script is in parent directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		scriptPath = filepath.Join(cwd, "..", "direct_api_post.py")
		logger.Printf("Using default Python script path: %s", scriptPath)
	}

	// Verify the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Fatalf("Python script not found at: %s\nPlease configure your config.py or use -python-script flag", scriptPath)
	}

	// Create bot configuration
	config := &bot.Config{
		TelegramToken:    *token,
		PythonScriptPath: scriptPath,
		DBPath:           *dbPath,
	}

	// Initialize bot
	logger.Println("Initializing bot...")
	logger.Printf("Python script: %s", scriptPath)
	logger.Printf("Database: %s", *dbPath)
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
	logger.Println("Make sure your config.py is configured with the correct BASE_URL and AGENCY_ID")
	<-sigChan

	// Graceful shutdown
	logger.Println("Shutting down...")
	poller.Stop()
	b.Stop()
	logger.Println("Goodbye!")
}
