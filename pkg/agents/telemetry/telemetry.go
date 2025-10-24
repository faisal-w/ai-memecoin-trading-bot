package telemetry

import (
	"log"
	"sync"
	"time"
)

// Metrics holds telemetry metrics
type Metrics struct {
	// Scanning metrics
	TokensScanned      int64
	TokensFound        int64
	
	// Filtering metrics
	TokensFiltered     int64
	TokensDropped      int64
	
	// Safety metrics
	SafetyChecks       int64
	HoneypotDetected   int64
	SafeTokens         int64
	
	// Strategy metrics
	Evaluations        int64
	CandidatesListed   int64
	TradesExecuted     int64
	
	// Execution metrics
	ExecutionSuccess   int64
	ExecutionFailed    int64
	SimulationFailed   int64
	
	// Financial metrics
	TotalInvested      float64
	TotalProfit        float64
	TotalLoss          float64
	
	// Performance metrics
	AvgDecisionLatency time.Duration
	AvgExecutionTime   time.Duration
	
	mu sync.RWMutex
}

// TelemetryAgent handles metrics and monitoring
type TelemetryAgent struct {
	metrics *Metrics
}

// NewTelemetryAgent creates a new telemetry agent
func NewTelemetryAgent() *TelemetryAgent {
	return &TelemetryAgent{
		metrics: &Metrics{},
	}
}

// RecordTokenScanned increments tokens scanned counter
func (t *TelemetryAgent) RecordTokenScanned() {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.TokensScanned++
}

// RecordTokenFound increments tokens found counter
func (t *TelemetryAgent) RecordTokenFound() {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.TokensFound++
}

// RecordTokenFiltered increments filtered counter
func (t *TelemetryAgent) RecordTokenFiltered(dropped bool) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.TokensFiltered++
	if dropped {
		t.metrics.TokensDropped++
	}
}

// RecordSafetyCheck records safety check result
func (t *TelemetryAgent) RecordSafetyCheck(isHoneypot, isSafe bool) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.SafetyChecks++
	if isHoneypot {
		t.metrics.HoneypotDetected++
	}
	if isSafe {
		t.metrics.SafeTokens++
	}
}

// RecordEvaluation increments evaluation counter
func (t *TelemetryAgent) RecordEvaluation() {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.Evaluations++
}

// RecordCandidateListed increments candidate listed counter
func (t *TelemetryAgent) RecordCandidateListed() {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.CandidatesListed++
}

// RecordExecution records trade execution result
func (t *TelemetryAgent) RecordExecution(success bool, amount float64) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.TradesExecuted++
	if success {
		t.metrics.ExecutionSuccess++
		t.metrics.TotalInvested += amount
	} else {
		t.metrics.ExecutionFailed++
	}
}

// RecordSimulationFailure increments simulation failure counter
func (t *TelemetryAgent) RecordSimulationFailure() {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	t.metrics.SimulationFailed++
}

// RecordProfit records profit/loss
func (t *TelemetryAgent) RecordProfit(profitLoss float64) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	if profitLoss > 0 {
		t.metrics.TotalProfit += profitLoss
	} else {
		t.metrics.TotalLoss += -profitLoss
	}
}

// RecordDecisionLatency records decision processing time
func (t *TelemetryAgent) RecordDecisionLatency(duration time.Duration) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	// Simple moving average
	t.metrics.AvgDecisionLatency = (t.metrics.AvgDecisionLatency + duration) / 2
}

// RecordExecutionTime records execution time
func (t *TelemetryAgent) RecordExecutionTime(duration time.Duration) {
	t.metrics.mu.Lock()
	defer t.metrics.mu.Unlock()
	// Simple moving average
	t.metrics.AvgExecutionTime = (t.metrics.AvgExecutionTime + duration) / 2
}

// GetMetrics returns a snapshot of current metrics
func (t *TelemetryAgent) GetMetrics() *Metrics {
	t.metrics.mu.RLock()
	defer t.metrics.mu.RUnlock()
	
	// Return a copy
	return &Metrics{
		TokensScanned:      t.metrics.TokensScanned,
		TokensFound:        t.metrics.TokensFound,
		TokensFiltered:     t.metrics.TokensFiltered,
		TokensDropped:      t.metrics.TokensDropped,
		SafetyChecks:       t.metrics.SafetyChecks,
		HoneypotDetected:   t.metrics.HoneypotDetected,
		SafeTokens:         t.metrics.SafeTokens,
		Evaluations:        t.metrics.Evaluations,
		CandidatesListed:   t.metrics.CandidatesListed,
		TradesExecuted:     t.metrics.TradesExecuted,
		ExecutionSuccess:   t.metrics.ExecutionSuccess,
		ExecutionFailed:    t.metrics.ExecutionFailed,
		SimulationFailed:   t.metrics.SimulationFailed,
		TotalInvested:      t.metrics.TotalInvested,
		TotalProfit:        t.metrics.TotalProfit,
		TotalLoss:          t.metrics.TotalLoss,
		AvgDecisionLatency: t.metrics.AvgDecisionLatency,
		AvgExecutionTime:   t.metrics.AvgExecutionTime,
	}
}

// LogMetrics logs current metrics
func (t *TelemetryAgent) LogMetrics() {
	metrics := t.GetMetrics()
	
	log.Println("=== Telemetry Metrics ===")
	log.Printf("Tokens Scanned: %d, Found: %d\n", metrics.TokensScanned, metrics.TokensFound)
	log.Printf("Tokens Filtered: %d, Dropped: %d\n", metrics.TokensFiltered, metrics.TokensDropped)
	log.Printf("Safety Checks: %d, Honeypots: %d, Safe: %d\n", 
		metrics.SafetyChecks, metrics.HoneypotDetected, metrics.SafeTokens)
	log.Printf("Evaluations: %d, Candidates: %d\n", metrics.Evaluations, metrics.CandidatesListed)
	log.Printf("Executions: %d (Success: %d, Failed: %d)\n", 
		metrics.TradesExecuted, metrics.ExecutionSuccess, metrics.ExecutionFailed)
	log.Printf("Financial: Invested: $%.2f, Profit: $%.2f, Loss: $%.2f\n",
		metrics.TotalInvested, metrics.TotalProfit, metrics.TotalLoss)
	log.Printf("Performance: Avg Decision: %v, Avg Execution: %v\n",
		metrics.AvgDecisionLatency, metrics.AvgExecutionTime)
	log.Println("========================")
}

// StartPeriodicLogging starts periodic metrics logging
func (t *TelemetryAgent) StartPeriodicLogging(interval time.Duration) chan struct{} {
	stopChan := make(chan struct{})
	
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				t.LogMetrics()
			case <-stopChan:
				return
			}
		}
	}()
	
	return stopChan
}
