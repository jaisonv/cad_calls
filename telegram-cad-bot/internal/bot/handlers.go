package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaisonv/telegram-cad-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) syncUserProfile(c tele.Context) {
	s := c.Sender()
	if s == nil {
		return
	}
	_ = b.db.UpsertUserProfile(s.ID, s.Username, s.FirstName, s.LastName)
}

// HandleStart handles the /start command
func (b *Bot) HandleStart(c tele.Context) error {
	b.syncUserProfile(c)
	msg := `👋 Welcome to CAD Calls Monitor Bot!

This bot monitors active CAD (Computer Aided Dispatch) calls and alerts you about incidents on streets you're watching.

📍 *Commands:*
/add <street_name> - Add a street to your watch list
/remove <street_name> - Remove a street from watch list
/list - Show all streets you're monitoring
/clear - Remove all monitored streets
/status - Show bot status and settings
/interval <minutes> - Set check interval (1-60 minutes)
/check - Manually check for new calls now
/help - Show this help message

*Example:*
/add Main Street
/add Oak Avenue

The bot will check for new calls every few minutes and notify you immediately when a call matches your watched streets.`

	return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// HandleHelp handles the /help command
func (b *Bot) HandleHelp(c tele.Context) error {
	return b.HandleStart(c)
}

// HandleAdd handles the /add command
func (b *Bot) HandleAdd(c tele.Context) error {
	b.syncUserProfile(c)
	args := strings.TrimSpace(c.Text())
	// Remove the command part
	streetName := strings.TrimPrefix(args, "/add")
	streetName = strings.TrimSpace(streetName)

	if streetName == "" {
		return c.Send("❌ Please provide a street name.\n\nExample: /add Main Street")
	}

	userID := c.Sender().ID
	if err := b.db.AddMonitoredStreet(userID, streetName); err != nil {
		b.logger.Printf("Error adding street for user %d: %v", userID, err)
		return c.Send("❌ Failed to add street. Please try again.")
	}

	return c.Send(fmt.Sprintf("✅ Added '%s' to your watch list!", streetName))
}

// HandleRemove handles the /remove command
func (b *Bot) HandleRemove(c tele.Context) error {
	b.syncUserProfile(c)
	args := strings.TrimSpace(c.Text())
	streetName := strings.TrimPrefix(args, "/remove")
	streetName = strings.TrimSpace(streetName)

	if streetName == "" {
		return c.Send("❌ Please provide a street name.\n\nExample: /remove Main Street")
	}

	userID := c.Sender().ID
	if err := b.db.RemoveMonitoredStreet(userID, streetName); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Send(fmt.Sprintf("❌ '%s' is not in your watch list.", streetName))
		}
		b.logger.Printf("Error removing street for user %d: %v", userID, err)
		return c.Send("❌ Failed to remove street. Please try again.")
	}

	return c.Send(fmt.Sprintf("✅ Removed '%s' from your watch list.", streetName))
}

// HandleList handles the /list command
func (b *Bot) HandleList(c tele.Context) error {
	b.syncUserProfile(c)
	userID := c.Sender().ID
	streets, err := b.db.GetMonitoredStreets(userID)
	if err != nil {
		b.logger.Printf("Error getting streets for user %d: %v", userID, err)
		return c.Send("❌ Failed to retrieve your watch list.")
	}

	if len(streets) == 0 {
		return c.Send("📋 Your watch list is empty.\n\nUse /add <street_name> to start monitoring streets.")
	}

	msg := fmt.Sprintf("📋 *Your Watch List* (%d streets):\n\n", len(streets))
	for i, street := range streets {
		msg += fmt.Sprintf("%d. %s\n", i+1, street)
	}

	return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// HandleClear handles the /clear command
func (b *Bot) HandleClear(c tele.Context) error {
	b.syncUserProfile(c)
	userID := c.Sender().ID
	if err := b.db.ClearMonitoredStreets(userID); err != nil {
		b.logger.Printf("Error clearing streets for user %d: %v", userID, err)
		return c.Send("❌ Failed to clear your watch list.")
	}

	return c.Send("✅ Your watch list has been cleared.")
}

// HandleStatus handles the /status command
func (b *Bot) HandleStatus(c tele.Context) error {
	b.syncUserProfile(c)
	userID := c.Sender().ID
	user, err := b.db.GetOrCreateUser(userID)
	if err != nil {
		b.logger.Printf("Error getting user %d: %v", userID, err)
		return c.Send("❌ Failed to retrieve your status.")
	}

	streets, err := b.db.GetMonitoredStreets(userID)
	if err != nil {
		b.logger.Printf("Error getting streets for user %d: %v", userID, err)
		return c.Send("❌ Failed to retrieve your status.")
	}

	msg := fmt.Sprintf(`📊 *Bot Status*

👤 User ID: %d
📍 Monitored Streets: %d
⏱ Check Interval: %d minutes
📅 Member Since: %s

🐍 *CAD Source:*
Python script (uses config.py settings)

Use /interval <minutes> to change check frequency (1-60 min).`,
		userID,
		len(streets),
		user.CheckInterval,
		user.CreatedAt.Format("2006-01-02"),
	)

	return c.Send(msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// HandleInterval handles the /interval command
func (b *Bot) HandleInterval(c tele.Context) error {
	b.syncUserProfile(c)
	args := strings.TrimSpace(c.Text())
	intervalStr := strings.TrimPrefix(args, "/interval")
	intervalStr = strings.TrimSpace(intervalStr)

	if intervalStr == "" {
		return c.Send("❌ Please provide an interval in minutes.\n\nExample: /interval 10")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval < 1 || interval > 60 {
		return c.Send("❌ Invalid interval. Please provide a number between 1 and 60 minutes.")
	}

	userID := c.Sender().ID
	if err := b.db.SetCheckInterval(userID, interval); err != nil {
		b.logger.Printf("Error setting interval for user %d: %v", userID, err)
		return c.Send("❌ Failed to set check interval.")
	}

	return c.Send(fmt.Sprintf("✅ Check interval set to %d minutes.", interval))
}

// HandleCheck handles the /check command (manual check)
func (b *Bot) HandleCheck(c tele.Context) error {
	b.syncUserProfile(c)
	userID := c.Sender().ID

	// Check if user has any monitored streets
	streets, err := b.db.GetMonitoredStreets(userID)
	if err != nil {
		return c.Send("❌ Failed to check for calls.")
	}

	if len(streets) == 0 {
		return c.Send("❌ You don't have any streets in your watch list.\n\nUse /add <street_name> to start monitoring.")
	}

	// Send "checking" message
	msg, _ := b.bot.Send(c.Recipient(), "🔍 Checking for new calls...")

	// Trigger a check for this user
	newCalls, err := b.CheckUserForNewCalls(userID)
	if err != nil {
		b.logger.Printf("Error checking calls for user %d: %v", userID, err)
		b.bot.Edit(msg, "❌ Failed to check for calls. Please try again later.")
		return nil
	}

	if len(newCalls) == 0 {
		b.bot.Edit(msg, "✅ No new calls found on your monitored streets.")
		return nil
	}

	b.bot.Edit(msg, fmt.Sprintf("✅ Found %d new call(s)! Sending details...", len(newCalls)))
	return nil
}

func (b *Bot) HandleUsers(c tele.Context) error {
	b.syncUserProfile(c)
	requesterID := c.Sender().ID
	if !b.isAdmin(requesterID) {
		return c.Send("❌ You are not authorized to use this command.")
	}

	users, err := b.db.GetAllUsersMonitoring()
	if err != nil {
		b.logger.Printf("Error getting users monitoring list: %v", err)
		return c.Send("❌ Failed to retrieve users.")
	}

	if len(users) == 0 {
		return c.Send("No users found.")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👥 *Users Monitoring Streets* (%d users)\n\n", len(users)))
	for _, u := range users {
		display := strings.TrimSpace(strings.TrimSpace(u.FirstName + " " + u.LastName))
		if display == "" {
			display = "(no name)"
		}
		sb.WriteString(fmt.Sprintf("• `%d`", u.TelegramID))
		if u.Username != "" {
			sb.WriteString(fmt.Sprintf(" @%s", u.Username))
		}
		sb.WriteString(fmt.Sprintf(" - %s\n", display))
		sb.WriteString(fmt.Sprintf("  streets (%d): ", u.StreetCount))
		if len(u.Streets) == 0 {
			sb.WriteString("none\n")
		} else {
			sb.WriteString(strings.Join(u.Streets, ", "))
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("  interval: %d min\n\n", u.CheckInterval))
	}

	return c.Send(sb.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// FormatCallMessage formats a CAD call into a readable message
func FormatCallMessage(call *storage.AlertCall) string {
	msg := fmt.Sprintf(`🚨 *New Call*

📍 *Address:* %s
🏷 *Type:* %s
📋 *Nature:* %s
🕐 *Time:* %s
🆔 *Incident:* %s`,
		call.Address,
		call.CallType,
		call.Nature,
		call.StartTime.Format("2006-01-02 15:04"),
		call.IncidentID,
	)

	if call.Agency != "" {
		msg += fmt.Sprintf("\n🏛 *Agency:* %s", call.Agency)
	}

	return msg
}
