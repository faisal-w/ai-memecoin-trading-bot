package risk

import (
	"log"
	"sync"
	"time"

	"github.com/mumugogoing/meme_bot/pkg/config"
	"github.com/mumugogoing/meme_bot/pkg/models"
)

// RiskManagerAgent manages risk controls and circuit breakers
type RiskManagerAgent struct {
	config  *config.Config
	control *models.RiskControl
	mu      sync.RWMutex
}

// NewRiskManagerAgent creates a new risk manager
func NewRiskManagerAgent(cfg *config.Config) *RiskManagerAgent {
	return &RiskManagerAgent{
		config: cfg,
		control: &models.RiskControl{
			SinglePositionPct: cfg.SinglePositionPct,
			TotalExposurePct:  cfg.TotalExposurePct,
			DailyLossLimit:    cfg.DailyLossLimit,
			CurrentExposure:   0,
			DailyLoss:         0,
			TradingHalted:     false,
			LastResetTime:     time.Now(),
		},
	}
}

// CanExecute checks if a trade can be executed based on risk controls
func (r *RiskManagerAgent) CanExecute(decision *models.StrategyDecision) (bool, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Check if trading is halted
	if r.control.TradingHalted {
		return false, "trading_halted"
	}
	
	// Check single position limit
	maxSinglePosition := r.config.AccountBalance * r.config.SinglePositionPct
	if decision.SuggestedAmountUSD > maxSinglePosition {
		log.Printf("RiskManager: Trade rejected - exceeds single position limit (%.2f > %.2f)\n",
			decision.SuggestedAmountUSD, maxSinglePosition)
		return false, "exceeds_single_position_limit"
	}
	
	// Check total exposure limit
	maxTotalExposure := r.config.AccountBalance * r.config.TotalExposurePct
	if r.control.CurrentExposure+decision.SuggestedAmountUSD > maxTotalExposure {
		log.Printf("RiskManager: Trade rejected - exceeds total exposure limit (%.2f > %.2f)\n",
			r.control.CurrentExposure+decision.SuggestedAmountUSD, maxTotalExposure)
		return false, "exceeds_total_exposure_limit"
	}
	
	// Check daily loss limit
	if r.control.DailyLoss >= r.config.DailyLossLimit {
		log.Printf("RiskManager: Trade rejected - daily loss limit reached (%.2f)\n",
			r.control.DailyLoss)
		r.haltTrading()
		return false, "daily_loss_limit_reached"
	}
	
	log.Printf("RiskManager: Trade approved for %s (%.2f USD)\n",
		decision.TokenAddress, decision.SuggestedAmountUSD)
	
	return true, ""
}

// RecordExecution records a trade execution
func (r *RiskManagerAgent) RecordExecution(result *models.ExecutionResult) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if result.Status == "confirmed" {
		r.control.CurrentExposure += result.AmountUSD
		log.Printf("RiskManager: Recorded execution - Current exposure: %.2f USD\n",
			r.control.CurrentExposure)
	}
}

// RecordProfit records profit/loss from a closed position
func (r *RiskManagerAgent) RecordProfit(tokenAddress string, profitLoss float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if profitLoss < 0 {
		r.control.DailyLoss += -profitLoss
		log.Printf("RiskManager: Recorded loss of %.2f - Daily loss: %.2f\n",
			-profitLoss, r.control.DailyLoss)
		
		// Check if daily loss limit reached
		if r.control.DailyLoss >= r.config.DailyLossLimit {
			r.haltTrading()
		}
	} else {
		log.Printf("RiskManager: Recorded profit of %.2f\n", profitLoss)
	}
}

// ReleaseExposure releases exposure when a position is closed
func (r *RiskManagerAgent) ReleaseExposure(amount float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.control.CurrentExposure -= amount
	if r.control.CurrentExposure < 0 {
		r.control.CurrentExposure = 0
	}
	
	log.Printf("RiskManager: Released exposure - Current exposure: %.2f USD\n",
		r.control.CurrentExposure)
}

// haltTrading triggers the circuit breaker
func (r *RiskManagerAgent) haltTrading() {
	r.control.TradingHalted = true
	log.Println("RiskManager: CIRCUIT BREAKER TRIGGERED - Trading halted!")
}

// ResumeTrading resumes trading (manual override)
func (r *RiskManagerAgent) ResumeTrading() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.control.TradingHalted = false
	log.Println("RiskManager: Trading resumed")
}

// ResetDaily resets daily counters (should be called at start of each day)
func (r *RiskManagerAgent) ResetDaily() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.control.DailyLoss = 0
	r.control.LastResetTime = time.Now()
	log.Println("RiskManager: Daily counters reset")
}

// CheckDailyReset checks if daily reset is needed
func (r *RiskManagerAgent) CheckDailyReset() {
	r.mu.RLock()
	lastReset := r.control.LastResetTime
	r.mu.RUnlock()
	
	// Reset if it's a new day
	if time.Since(lastReset) > 24*time.Hour {
		r.ResetDaily()
	}
}

// GetStatus returns current risk control status
func (r *RiskManagerAgent) GetStatus() *models.RiskControl {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return a copy
	return &models.RiskControl{
		SinglePositionPct: r.control.SinglePositionPct,
		TotalExposurePct:  r.control.TotalExposurePct,
		DailyLossLimit:    r.control.DailyLossLimit,
		CurrentExposure:   r.control.CurrentExposure,
		DailyLoss:         r.control.DailyLoss,
		TradingHalted:     r.control.TradingHalted,
		LastResetTime:     r.control.LastResetTime,
	}
}
