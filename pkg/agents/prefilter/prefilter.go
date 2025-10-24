package prefilter

import (
	"log"
	"strings"

	"github.com/mumugogoing/meme_bot/pkg/config"
	"github.com/mumugogoing/meme_bot/pkg/models"
)

// PreFilterAgent performs basic filtering on discovered tokens
type PreFilterAgent struct {
	config *config.Config
}

// NewPreFilterAgent creates a new pre-filter agent
func NewPreFilterAgent(cfg *config.Config) *PreFilterAgent {
	return &PreFilterAgent{
		config: cfg,
	}
}

// Filter applies pre-filtering rules to a token
func (p *PreFilterAgent) Filter(token models.TokenFound) models.PreFilteredToken {
	result := models.PreFilteredToken{
		Token:    token,
		Priority: "medium",
		Dropped:  false,
		Reasons:  []string{},
	}
	
	// Check blacklisted tokens
	if p.isBlacklistedToken(token.TokenAddress) {
		result.Dropped = true
		result.Reasons = append(result.Reasons, "token_blacklisted")
		log.Printf("PreFilterAgent: Token %s dropped - blacklisted\n", token.TokenAddress)
		return result
	}
	
	// Check blacklisted creators
	if p.isBlacklistedCreator(token.CreatorAddress) {
		result.Dropped = true
		result.Reasons = append(result.Reasons, "creator_blacklisted")
		log.Printf("PreFilterAgent: Token %s dropped - creator blacklisted\n", token.TokenAddress)
		return result
	}
	
	// Check whitelisted tokens (high priority)
	if p.isWhitelistedToken(token.TokenAddress) {
		result.Priority = "high"
		result.Reasons = append(result.Reasons, "token_whitelisted")
		log.Printf("PreFilterAgent: Token %s marked high priority - whitelisted\n", token.TokenAddress)
		return result
	}
	
	// Check minimum liquidity
	totalLiquidity := token.InitialLiquidity.ReserveNative
	if totalLiquidity < p.config.MinLiquidity {
		result.Priority = "low"
		result.Reasons = append(result.Reasons, "low_initial_liquidity")
		log.Printf("PreFilterAgent: Token %s marked low priority - low liquidity (%.2f)\n", 
			token.TokenAddress, totalLiquidity)
	}
	
	// Check for suspicious patterns in metadata
	if p.hasSuspiciousMetadata(token) {
		result.Priority = "low"
		result.Reasons = append(result.Reasons, "suspicious_metadata")
		log.Printf("PreFilterAgent: Token %s marked low priority - suspicious metadata\n", token.TokenAddress)
	}
	
	// Check for very high initial liquidity (potential whale)
	if totalLiquidity > 100000 {
		result.Priority = "high"
		result.Reasons = append(result.Reasons, "high_initial_liquidity")
		log.Printf("PreFilterAgent: Token %s marked high priority - high liquidity (%.2f)\n", 
			token.TokenAddress, totalLiquidity)
	}
	
	return result
}

// isBlacklistedToken checks if token is in blacklist
func (p *PreFilterAgent) isBlacklistedToken(address string) bool {
	for _, blacklisted := range p.config.BlacklistedTokens {
		if strings.EqualFold(address, blacklisted) {
			return true
		}
	}
	return false
}

// isBlacklistedCreator checks if creator is in blacklist
func (p *PreFilterAgent) isBlacklistedCreator(address string) bool {
	for _, blacklisted := range p.config.BlacklistedCreators {
		if strings.EqualFold(address, blacklisted) {
			return true
		}
	}
	return false
}

// isWhitelistedToken checks if token is in whitelist
func (p *PreFilterAgent) isWhitelistedToken(address string) bool {
	for _, whitelisted := range p.config.WhitelistedTokens {
		if strings.EqualFold(address, whitelisted) {
			return true
		}
	}
	return false
}

// hasSuspiciousMetadata checks for suspicious patterns in token metadata
func (p *PreFilterAgent) hasSuspiciousMetadata(token models.TokenFound) bool {
	// Check for common scam patterns in metadata
	suspiciousWords := []string{
		"test", "scam", "rug", "fake", "honeypot",
		"xxx", "pump", "dump", "bot",
	}
	
	for key, value := range token.Metadata {
		lowerKey := strings.ToLower(key)
		lowerValue := strings.ToLower(value)
		
		for _, word := range suspiciousWords {
			if strings.Contains(lowerKey, word) || strings.Contains(lowerValue, word) {
				return true
			}
		}
	}
	
	return false
}
