package strategy

import (
	"log"
	"math"
	"time"

	"github.com/mumugogoing/meme_bot/pkg/config"
	"github.com/mumugogoing/meme_bot/pkg/models"
)

// StrategyEvaluatorAgent evaluates trading opportunities
type StrategyEvaluatorAgent struct {
	config *config.Config
}

// NewStrategyEvaluatorAgent creates a new strategy evaluator
func NewStrategyEvaluatorAgent(cfg *config.Config) *StrategyEvaluatorAgent {
	return &StrategyEvaluatorAgent{
		config: cfg,
	}
}

// Evaluate performs comprehensive strategy evaluation
func (s *StrategyEvaluatorAgent) Evaluate(
	safety *models.SafetyReport,
	offchain *models.OffChainMetrics,
	token models.PreFilteredToken,
) (*models.StrategyDecision, error) {
	log.Printf("StrategyEvaluatorAgent: Evaluating token %s\n", token.Token.TokenAddress)
	
	decision := &models.StrategyDecision{
		TokenAddress:   token.Token.TokenAddress,
		Chain:          token.Token.Chain,
		EvaluatedAt:    time.Now(),
		Rationale:      []string{},
	}
	
	// Calculate win probability
	decision.WinProbability = s.calculateWinProbability(safety, offchain, token)
	
	// Calculate expected ROI
	decision.ExpectedROI, decision.ExpectedROIStd = s.calculateExpectedROI(safety, offchain, token)
	
	// Determine confidence level
	decision.Confidence = s.determineConfidence(safety, offchain, decision.WinProbability)
	
	// Determine action
	decision.Action = s.determineAction(decision)
	
	// Calculate suggested position size
	decision.SuggestedAmountUSD = s.calculatePositionSize(decision)
	
	// Set risk parameters
	decision.StopLossPct = s.calculateStopLoss(decision)
	decision.TakeProfitPct = s.calculateTakeProfit(decision)
	decision.TimeHorizonMinutes = s.calculateTimeHorizon(decision)
	
	log.Printf("StrategyEvaluatorAgent: Token %s - WinProb: %.2f, Action: %s, Confidence: %s\n",
		token.Token.TokenAddress, decision.WinProbability, decision.Action, decision.Confidence)
	
	return decision, nil
}

// calculateWinProbability calculates the probability of a profitable trade
func (s *StrategyEvaluatorAgent) calculateWinProbability(
	safety *models.SafetyReport,
	offchain *models.OffChainMetrics,
	token models.PreFilteredToken,
) float64 {
	baseProb := 0.5 // Start at 50%
	
	// Safety factors (most important)
	if safety.CanBuy && safety.CanSell {
		baseProb += 0.15
		token.Reasons = append(token.Reasons, "can_trade")
	} else {
		return 0.0 // Cannot trade = 0% win probability
	}
	
	if safety.HoneypotScore < 0.1 {
		baseProb += 0.10
		token.Reasons = append(token.Reasons, "low_honeypot_score")
	} else if safety.HoneypotScore > s.config.MaxHoneypotScore {
		baseProb -= 0.20
		token.Reasons = append(token.Reasons, "high_honeypot_score")
	}
	
	if safety.LiquidityLocked {
		baseProb += 0.08
		token.Reasons = append(token.Reasons, "liquidity_locked")
	}
	
	if safety.OwnerControls.Renounced {
		baseProb += 0.07
		token.Reasons = append(token.Reasons, "owner_renounced")
	}
	
	if !safety.OwnerControls.HasBlacklist && !safety.OwnerControls.HasTransferHook {
		baseProb += 0.05
		token.Reasons = append(token.Reasons, "no_transfer_restrictions")
	}
	
	// Volume factors
	if offchain.Volume24hDEX >= s.config.MinVolumeDEX {
		baseProb += 0.10
		token.Reasons = append(token.Reasons, "good_dex_volume")
	}
	
	// Social factors
	totalMentions := 0
	for _, count := range offchain.SocialMentions {
		totalMentions += count
	}
	if totalMentions > 50 {
		baseProb += 0.08
		token.Reasons = append(token.Reasons, "social_activity")
	}
	
	// Velocity factor
	if offchain.Velocity == "rising" {
		baseProb += 0.07
		token.Reasons = append(token.Reasons, "rising_velocity")
	} else if offchain.Velocity == "falling" {
		baseProb -= 0.10
	}
	
	// Liquidity concentration
	liquidityRatio := token.Token.InitialLiquidity.ReserveNative / (token.Token.InitialLiquidity.ReserveNative + token.Token.InitialLiquidity.ReserveToken)
	if liquidityRatio < 0.3 || liquidityRatio > 0.7 {
		token.Reasons = append(token.Reasons, "liquidity_imbalance")
		baseProb -= 0.05
	}
	
	// Priority adjustment
	if token.Priority == "high" {
		baseProb += 0.05
	} else if token.Priority == "low" {
		baseProb -= 0.05
	}
	
	// Clamp between 0 and 1
	if baseProb > 1.0 {
		baseProb = 1.0
	}
	if baseProb < 0.0 {
		baseProb = 0.0
	}
	
	return baseProb
}

// calculateExpectedROI calculates expected return and standard deviation
func (s *StrategyEvaluatorAgent) calculateExpectedROI(
	safety *models.SafetyReport,
	offchain *models.OffChainMetrics,
	token models.PreFilteredToken,
) (float64, float64) {
	// Base expected return for meme coins
	baseROI := 0.15
	
	// Adjust based on volume and liquidity
	if offchain.Volume24hDEX > s.config.MinVolumeDEX*2 {
		baseROI += 0.10
	}
	
	if token.Token.InitialLiquidity.ReserveNative > s.config.MinLiquidity*2 {
		baseROI += 0.08
	}
	
	// Adjust for velocity
	if offchain.Velocity == "rising" {
		baseROI += 0.12
	}
	
	// Standard deviation (volatility measure)
	stdDev := 0.25 // High volatility for meme coins
	
	return baseROI, stdDev
}

// determineConfidence determines confidence level
func (s *StrategyEvaluatorAgent) determineConfidence(
	safety *models.SafetyReport,
	offchain *models.OffChainMetrics,
	winProb float64,
) string {
	// High confidence requires multiple factors
	if winProb >= 0.85 && 
	   safety.HoneypotScore < 0.1 && 
	   offchain.Volume24hDEX > s.config.MinVolumeDEX {
		return "high"
	}
	
	if winProb >= 0.70 && safety.HoneypotScore < 0.2 {
		return "medium"
	}
	
	return "low"
}

// determineAction determines the recommended action
func (s *StrategyEvaluatorAgent) determineAction(decision *models.StrategyDecision) string {
	// Check if meets threshold for listing/buying
	if decision.WinProbability >= s.config.WinProbabilityThreshold && 
	   decision.Confidence != "low" {
		if s.config.AutoExecute {
			return "buy"
		}
		return "list"
	}
	
	if decision.WinProbability >= 0.60 {
		return "monitor"
	}
	
	return "skip"
}

// calculatePositionSize calculates suggested position size
func (s *StrategyEvaluatorAgent) calculatePositionSize(decision *models.StrategyDecision) float64 {
	// Kelly Criterion simplified: f = (p * b - q) / b
	// where p = win probability, q = 1-p, b = odds (ROI)
	
	maxPosition := s.config.AccountBalance * s.config.SinglePositionPct
	
	// Adjust based on confidence
	multiplier := 1.0
	switch decision.Confidence {
	case "high":
		multiplier = 1.0
	case "medium":
		multiplier = 0.7
	case "low":
		multiplier = 0.4
	}
	
	suggestedAmount := maxPosition * multiplier
	
	// Minimum viable position
	if suggestedAmount < 100 {
		suggestedAmount = 100
	}
	
	return math.Round(suggestedAmount*100) / 100
}

// calculateStopLoss calculates stop loss percentage
func (s *StrategyEvaluatorAgent) calculateStopLoss(decision *models.StrategyDecision) float64 {
	// Tighter stop loss for lower confidence
	switch decision.Confidence {
	case "high":
		return 0.20
	case "medium":
		return 0.15
	case "low":
		return 0.10
	}
	return 0.15
}

// calculateTakeProfit calculates take profit percentage
func (s *StrategyEvaluatorAgent) calculateTakeProfit(decision *models.StrategyDecision) float64 {
	// Higher take profit for higher expected ROI
	baseTP := decision.ExpectedROI * 1.5
	
	if baseTP < 0.20 {
		baseTP = 0.20
	}
	if baseTP > 1.00 {
		baseTP = 1.00
	}
	
	return math.Round(baseTP*100) / 100
}

// calculateTimeHorizon calculates recommended holding period
func (s *StrategyEvaluatorAgent) calculateTimeHorizon(decision *models.StrategyDecision) int {
	// Meme coins are typically short-term trades
	// Base on confidence and expected ROI
	
	if decision.Confidence == "high" {
		return 60 // 1 hour
	}
	if decision.Confidence == "medium" {
		return 30 // 30 minutes
	}
	return 15 // 15 minutes for low confidence
}
