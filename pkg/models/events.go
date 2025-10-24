package models

import "time"

// Chain represents supported blockchain networks
type Chain string

const (
	ChainSolana Chain = "solana"
	ChainBase   Chain = "base"
)

// TokenFound event from ChainScannerAgent
type TokenFound struct {
	Chain            Chain              `json:"chain"`
	TokenAddress     string             `json:"token_address"`
	FirstSeenTS      int64              `json:"first_seen_ts"`
	CreatorAddress   string             `json:"creator_address"`
	InitialLiquidity InitialLiquidity   `json:"initial_liquidity"`
	TxHash           string             `json:"tx_hash"`
	Metadata         map[string]string  `json:"metadata,omitempty"`
}

// InitialLiquidity details
type InitialLiquidity struct {
	Pair          string  `json:"pair"`
	ReserveToken  float64 `json:"reserve_token"`
	ReserveNative float64 `json:"reserve_native"`
}

// PreFilteredToken event after basic filtering
type PreFilteredToken struct {
	Token    TokenFound `json:"token"`
	Priority string     `json:"priority"` // "high", "medium", "low"
	Dropped  bool       `json:"dropped"`
	Reasons  []string   `json:"reasons,omitempty"`
}

// SafetyReport from OnChainSafetyAgent
type SafetyReport struct {
	TokenAddress      string               `json:"token_address"`
	Chain             Chain                `json:"chain"`
	CanBuy            bool                 `json:"can_buy"`
	CanSell           bool                 `json:"can_sell"`
	HoneypotScore     float64              `json:"honeypot_score"` // 0..1
	LiquidityLocked   bool                 `json:"liquidity_locked"`
	OwnerControls     OwnerControls        `json:"owner_controls"`
	SimulatedSell     SimulatedSellResult  `json:"simulated_sell_result"`
	Reasons           []string             `json:"reasons,omitempty"`
	EvaluatedAt       time.Time            `json:"evaluated_at"`
}

// OwnerControls details
type OwnerControls struct {
	Renounced       bool    `json:"renounced"`
	HasBlacklist    bool    `json:"has_blacklist"`
	MaxTxLimit      float64 `json:"max_tx_limit,omitempty"`
	TaxFee          float64 `json:"tax_fee,omitempty"`
	HasTransferHook bool    `json:"has_transfer_hook"`
}

// SimulatedSellResult details
type SimulatedSellResult struct {
	Success  bool    `json:"success"`
	Slippage float64 `json:"slippage"`
	GasUsed  uint64  `json:"gas_used,omitempty"`
	Error    string  `json:"error,omitempty"`
}

// OffChainMetrics from OffChainDataAgent
type OffChainMetrics struct {
	TokenAddress    string                 `json:"token_address"`
	Volume24hCEX    float64                `json:"24h_volume_cex"`
	Volume24hDEX    float64                `json:"24h_volume_dex"`
	SocialMentions  map[string]int         `json:"social_mentions"` // twitter, telegram, reddit
	Velocity        string                 `json:"velocity"`        // "rising", "stable", "falling"
	PriceOnCEX      float64                `json:"price_on_cex,omitempty"`
	PriceOnDEX      float64                `json:"price_on_dex,omitempty"`
	MarketCap       float64                `json:"market_cap,omitempty"`
	EvaluatedAt     time.Time              `json:"evaluated_at"`
}

// StrategyDecision from StrategyEvaluatorAgent
type StrategyDecision struct {
	TokenAddress        string    `json:"token_address"`
	Chain               Chain     `json:"chain"`
	WinProbability      float64   `json:"win_probability"`      // 0..1
	ExpectedROI         float64   `json:"expected_roi"`         // mean ROI
	ExpectedROIStd      float64   `json:"expected_roi_std"`     // std deviation
	Confidence          string    `json:"confidence"`           // "high", "medium", "low"
	Action              string    `json:"action"`               // "list", "buy", "monitor", "skip"
	SuggestedAmountUSD  float64   `json:"suggested_amount_usd"`
	StopLossPct         float64   `json:"stop_loss_pct"`
	TakeProfitPct       float64   `json:"take_profit_pct"`
	TimeHorizonMinutes  int       `json:"time_horizon_minutes"`
	EvaluatedAt         time.Time `json:"evaluated_at"`
	Rationale           []string  `json:"rationale,omitempty"`
}

// CandidateToken for listing queue
type CandidateToken struct {
	Token           TokenFound        `json:"token"`
	SafetyReport    SafetyReport      `json:"safety_report"`
	OffChainMetrics OffChainMetrics   `json:"offchain_metrics"`
	StrategyDecision StrategyDecision `json:"strategy_decision"`
	ListedAt        time.Time         `json:"listed_at"`
	Status          string            `json:"status"` // "pending", "approved", "rejected", "executed"
}

// ExecutionResult from ExecutionAgent
type ExecutionResult struct {
	TokenAddress   string    `json:"token_address"`
	Chain          Chain     `json:"chain"`
	TxHash         string    `json:"tx_hash"`
	Status         string    `json:"status"` // "pending", "confirmed", "failed"
	GasUsed        uint64    `json:"gas_used,omitempty"`
	SlippageActual float64   `json:"slippage_actual,omitempty"`
	AmountUSD      float64   `json:"amount_usd"`
	Timestamp      time.Time `json:"timestamp"`
	Error          string    `json:"error,omitempty"`
}

// RiskControl parameters and state
type RiskControl struct {
	SinglePositionPct  float64 `json:"single_position_pct"`  // max % of balance per trade
	TotalExposurePct   float64 `json:"total_exposure_pct"`   // max % of total balance exposed
	DailyLossLimit     float64 `json:"daily_loss_limit"`     // max daily loss in USD
	CurrentExposure    float64 `json:"current_exposure"`     // current total exposure
	DailyLoss          float64 `json:"daily_loss"`           // current daily loss
	TradingHalted      bool    `json:"trading_halted"`       // circuit breaker status
	LastResetTime      time.Time `json:"last_reset_time"`
}
