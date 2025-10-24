package listing

import (
	"log"
	"sync"
	"time"

	"github.com/mumugogoing/meme_bot/pkg/models"
)

// CandidateListingAgent manages the candidate token queue
type CandidateListingAgent struct {
	candidates map[string]*models.CandidateToken
	queue      chan *models.CandidateToken
	mu         sync.RWMutex
}

// NewCandidateListingAgent creates a new listing agent
func NewCandidateListingAgent() *CandidateListingAgent {
	return &CandidateListingAgent{
		candidates: make(map[string]*models.CandidateToken),
		queue:      make(chan *models.CandidateToken, 100),
	}
}

// AddCandidate adds a token to the candidate list
func (c *CandidateListingAgent) AddCandidate(
	token models.TokenFound,
	safety models.SafetyReport,
	offchain models.OffChainMetrics,
	decision models.StrategyDecision,
) *models.CandidateToken {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	candidate := &models.CandidateToken{
		Token:            token,
		SafetyReport:     safety,
		OffChainMetrics:  offchain,
		StrategyDecision: decision,
		ListedAt:         time.Now(),
		Status:           "pending",
	}
	
	c.candidates[token.TokenAddress] = candidate
	
	log.Printf("CandidateListingAgent: Added candidate %s - WinProb: %.2f, Action: %s\n",
		token.TokenAddress, decision.WinProbability, decision.Action)
	
	// Send to queue for execution consideration
	select {
	case c.queue <- candidate:
	default:
		log.Println("CandidateListingAgent: Warning - queue full, candidate not queued")
	}
	
	return candidate
}

// GetCandidate retrieves a candidate by token address
func (c *CandidateListingAgent) GetCandidate(tokenAddress string) (*models.CandidateToken, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	candidate, exists := c.candidates[tokenAddress]
	return candidate, exists
}

// GetAllCandidates returns all candidates
func (c *CandidateListingAgent) GetAllCandidates() []*models.CandidateToken {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	candidates := make([]*models.CandidateToken, 0, len(c.candidates))
	for _, candidate := range c.candidates {
		candidates = append(candidates, candidate)
	}
	
	return candidates
}

// GetPendingCandidates returns candidates with pending status
func (c *CandidateListingAgent) GetPendingCandidates() []*models.CandidateToken {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	candidates := make([]*models.CandidateToken, 0)
	for _, candidate := range c.candidates {
		if candidate.Status == "pending" {
			candidates = append(candidates, candidate)
		}
	}
	
	return candidates
}

// UpdateStatus updates the status of a candidate
func (c *CandidateListingAgent) UpdateStatus(tokenAddress, status string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if candidate, exists := c.candidates[tokenAddress]; exists {
		candidate.Status = status
		log.Printf("CandidateListingAgent: Updated %s status to %s\n", tokenAddress, status)
	}
}

// GetQueue returns the candidate queue channel
func (c *CandidateListingAgent) GetQueue() <-chan *models.CandidateToken {
	return c.queue
}

// GetCandidateCount returns the total number of candidates
func (c *CandidateListingAgent) GetCandidateCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.candidates)
}
