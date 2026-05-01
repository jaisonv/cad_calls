package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/jaisonv/telegram-cad-bot/internal/storage"
)

// UserChecker is an interface for checking calls for a user
type UserChecker interface {
	CheckUserForNewCalls(userID int64) ([]*storage.AlertCall, error)
	GetDB() *storage.DB
}

// Poller handles periodic checking for new CAD calls
type Poller struct {
	checker       UserChecker
	db            *storage.DB
	logger        *log.Logger
	stopChan      chan struct{}
	wg            sync.WaitGroup
	checkInterval time.Duration
	cleanupTicker *time.Ticker
}

// NewPoller creates a new poller instance
func NewPoller(checker UserChecker, db *storage.DB, checkInterval time.Duration, logger *log.Logger) *Poller {
	return &Poller{
		checker:       checker,
		db:            db,
		logger:        logger,
		stopChan:      make(chan struct{}),
		checkInterval: checkInterval,
		cleanupTicker: time.NewTicker(24 * time.Hour), // Daily cleanup
	}
}

// Start begins the polling process
func (p *Poller) Start() {
	p.logger.Println("Poller started")

	// Start the main polling loop
	p.wg.Add(1)
	go p.pollLoop()

	// Start the cleanup routine
	p.wg.Add(1)
	go p.cleanupLoop()
}

// Stop stops the polling process
func (p *Poller) Stop() {
	p.logger.Println("Stopping poller...")
	close(p.stopChan)
	p.cleanupTicker.Stop()
	p.wg.Wait()
	p.logger.Println("Poller stopped")
}

// pollLoop runs the main polling loop
func (p *Poller) pollLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.checkInterval)
	defer ticker.Stop()

	// Do an initial check immediately
	p.checkAllUsers()

	for {
		select {
		case <-ticker.C:
			p.checkAllUsers()
		case <-p.stopChan:
			return
		}
	}
}

// cleanupLoop periodically cleans up old seen calls
func (p *Poller) cleanupLoop() {
	defer p.wg.Done()

	for {
		select {
		case <-p.cleanupTicker.C:
			p.logger.Println("Running cleanup of old seen calls...")
			if err := p.db.CleanupOldSeenCalls(7); err != nil {
				p.logger.Printf("Error during cleanup: %v", err)
			} else {
				p.logger.Println("Cleanup completed")
			}
		case <-p.stopChan:
			return
		}
	}
}

// checkAllUsers checks for new calls for all active users
func (p *Poller) checkAllUsers() {
	users, err := p.db.GetAllActiveUsers()
	if err != nil {
		p.logger.Printf("Error getting active users: %v", err)
		return
	}

	if len(users) == 0 {
		p.logger.Println("No active users to check")
		return
	}

	p.logger.Printf("Checking calls for %d active users", len(users))

	// Check each user concurrently (with rate limiting)
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent checks
	var wg sync.WaitGroup

	for _, userID := range users {
		wg.Add(1)
		go func(uid int64) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			newCalls, err := p.checker.CheckUserForNewCalls(uid)
			if err != nil {
				p.logger.Printf("Error checking calls for user %d: %v", uid, err)
				return
			}

			if len(newCalls) > 0 {
				p.logger.Printf("Found %d new calls for user %d", len(newCalls), uid)
			}
		}(userID)
	}

	wg.Wait()
}
