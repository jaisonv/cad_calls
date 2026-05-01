package bot

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jaisonv/telegram-cad-bot/internal/cad"
	"github.com/jaisonv/telegram-cad-bot/internal/filter"
	"github.com/jaisonv/telegram-cad-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

// Config holds the bot configuration
type Config struct {
	TelegramToken    string
	PythonScriptPath string
	DBPath           string
	AdminUserIDs     map[int64]struct{}
}

// Bot represents the Telegram bot
type Bot struct {
	bot       *tele.Bot
	db        *storage.DB
	cadClient *cad.Client
	config    *Config
	logger    *log.Logger
}

// NewBot creates a new Telegram bot instance
func NewBot(config *Config, logger *log.Logger) (*Bot, error) {
	// Initialize database
	db, err := storage.NewDB(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize CAD client (uses Python script)
	cadClient := cad.NewClient(config.PythonScriptPath, "")

	// Initialize Telegram bot
	pref := tele.Settings{
		Token:  config.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	bot := &Bot{
		bot:       b,
		db:        db,
		cadClient: cadClient,
		config:    config,
		logger:    logger,
	}

	// Register handlers
	bot.registerHandlers()
	if err := bot.configureCommands(); err != nil {
		logger.Printf("Warning: failed to configure bot commands: %v", err)
	}

	return bot, nil
}

// registerHandlers registers all command handlers
func (b *Bot) registerHandlers() {
	b.bot.Handle("/start", b.HandleStart)
	b.bot.Handle("/help", b.HandleHelp)
	b.bot.Handle("/add", b.HandleAdd)
	b.bot.Handle("/remove", b.HandleRemove)
	b.bot.Handle("/list", b.HandleList)
	b.bot.Handle("/clear", b.HandleClear)
	b.bot.Handle("/status", b.HandleStatus)
	b.bot.Handle("/interval", b.HandleInterval)
	b.bot.Handle("/check", b.HandleCheck)
	b.bot.Handle("/users", b.HandleUsers)
}

func (b *Bot) isAdmin(userID int64) bool {
	_, ok := b.config.AdminUserIDs[userID]
	return ok
}

func (b *Bot) configureCommands() error {
	// Clear broad scopes first to avoid stale command visibility
	// from previous BotFather/API setups.
	_ = b.bot.DeleteCommands(&tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	_ = b.bot.DeleteCommands(&tele.CommandScope{Type: tele.CommandScopeAllGroupChats})
	_ = b.bot.DeleteCommands(&tele.CommandScope{Type: tele.CommandScopeAllChatAdmin})

	defaultCommands := []tele.Command{
		{Text: "start", Description: "Start the bot"},
		{Text: "help", Description: "Show help"},
		{Text: "add", Description: "Add street to watch"},
		{Text: "remove", Description: "Remove watched street"},
		{Text: "list", Description: "Show watched streets"},
		{Text: "clear", Description: "Clear watch list"},
		{Text: "status", Description: "Show bot status"},
		{Text: "interval", Description: "Set check interval"},
		{Text: "check", Description: "Check calls now"},
	}

	if err := b.bot.SetCommands(defaultCommands, &tele.CommandScope{Type: tele.CommandScopeDefault}); err != nil {
		return err
	}

	// Keep /users hidden from command suggestions.
	// Admins can still run it manually and it remains access-controlled in code.
	adminIDs := make([]int64, 0, len(b.config.AdminUserIDs))
	for id := range b.config.AdminUserIDs {
		adminIDs = append(adminIDs, id)
	}
	sort.Slice(adminIDs, func(i, j int) bool { return adminIDs[i] < adminIDs[j] })
	for _, adminID := range adminIDs {
		_ = b.bot.DeleteCommands(&tele.CommandScope{Type: tele.CommandScopeChat, ChatID: adminID})
	}

	return nil
}

// Start starts the bot
func (b *Bot) Start() {
	b.logger.Println("Bot is starting...")
	b.bot.Start()
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.logger.Println("Bot is stopping...")
	b.bot.Stop()
	b.db.Close()
}

// CheckUserForNewCalls checks for new calls for a specific user and sends alerts
func (b *Bot) CheckUserForNewCalls(userID int64) ([]*storage.AlertCall, error) {
	// Get user's monitored streets
	streets, err := b.db.GetMonitoredStreets(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitored streets: %w", err)
	}

	if len(streets) == 0 {
		return nil, nil
	}

	// Fetch active calls
	resp, err := b.cadClient.GetActiveCalls(50)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CAD calls: %w", err)
	}

	// Filter calls by monitored streets
	matchedCalls := filter.FilterCalls(resp.CADCalls, streets)

	// Check which calls are new
	var newCalls []*storage.AlertCall
	for _, call := range matchedCalls {
		seen, err := b.db.HasSeenCall(userID, call.IncidentID)
		if err != nil {
			b.logger.Printf("Error checking if call was seen: %v", err)
			continue
		}

		if !seen {
			// Mark as seen
			if err := b.db.MarkCallAsSeen(userID, call.IncidentID); err != nil {
				b.logger.Printf("Error marking call as seen: %v", err)
			}

			// Add to new calls list
			alertCall := &storage.AlertCall{
				IncidentID: call.IncidentID,
				CallType:   call.CallType,
				Nature:     call.Nature,
				Address:    call.Address,
				StartTime:  call.StartTime,
				Agency:     call.Agency,
			}
			newCalls = append(newCalls, alertCall)

			// Send alert
			b.SendAlert(userID, alertCall)
		}
	}

	return newCalls, nil
}

// SendAlert sends an alert message to a user
func (b *Bot) SendAlert(userID int64, call *storage.AlertCall) {
	recipient := &tele.User{ID: userID}
	msg := FormatCallMessage(call)

	_, err := b.bot.Send(recipient, msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	if err != nil {
		b.logger.Printf("Error sending alert to user %d: %v", userID, err)
	}
}

// GetDB returns the database instance
func (b *Bot) GetDB() *storage.DB {
	return b.db
}
