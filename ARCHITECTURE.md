# Trading Bot Architecture

## System Design

The meme coin trading bot uses a multi-agent architecture where each agent has a specific responsibility. Agents communicate via channels and shared data structures.

## Agents Overview

### 1. ChainScannerAgent
**Purpose**: Monitor blockchain networks for new token creation events

**Responsibilities**:
- Subscribe to Solana and Base RPC/WebSocket endpoints
- Monitor token factory contracts (Uniswap, Raydium, Orca)
- Parse transaction logs for new token mints
- Extract initial liquidity and creator information
- Emit `TokenFound` events

**Output**: `TokenFound` events via channel

### 2. PreFilterAgent
**Purpose**: Apply fast, simple rules to filter out obvious junk tokens

**Responsibilities**:
- Check blacklisted token addresses
- Check blacklisted creator addresses
- Check whitelisted token addresses
- Validate minimum initial liquidity
- Detect suspicious metadata patterns
- Assign priority (high/medium/low)

**Output**: `PreFilteredToken` with priority and drop status

### 3. OnChainSafetyAgent
**Purpose**: Perform comprehensive on-chain safety analysis (honeypot detection)

**Responsibilities**:
- Simulate buy transaction
- Simulate sell transaction
- Detect transfer restrictions (blacklist, hooks)
- Check owner controls (renounced, permissions)
- Analyze tax fees
- Verify liquidity lock status
- Calculate honeypot risk score (0-1)

**Output**: `SafetyReport` with can_buy/can_sell flags and score

**Key Metrics**:
- `can_buy`: Can tokens be purchased
- `can_sell`: Can tokens be sold (critical for honeypot detection)
- `honeypot_score`: Overall risk score (0 = safe, 1 = definite honeypot)
- `slippage`: Expected slippage on sell

### 4. OffChainDataAgent
**Purpose**: Gather trading volume and social media signals

**Responsibilities**:
- Query DEX aggregators for trading volume
- Query CEX APIs for listed tokens
- Fetch social media mentions (Twitter, Telegram, Reddit)
- Calculate price from DEX and CEX
- Determine velocity trend (rising/stable/falling)

**Output**: `OffChainMetrics` with volume, social, and price data

**Data Sources**:
- TheGraph (DEX data)
- CoinGecko API
- OKX API
- Twitter API
- Telegram/Reddit scrapers

### 5. StrategyEvaluatorAgent
**Purpose**: Calculate win probability and recommend trading actions

**Responsibilities**:
- Combine safety report and off-chain metrics
- Calculate win probability (0-1)
- Estimate expected ROI
- Determine confidence level (high/medium/low)
- Recommend action (buy/list/monitor/skip)
- Calculate position size
- Set stop-loss and take-profit levels

**Output**: `StrategyDecision` with win probability and action

**Algorithm**:
```
Base Win Probability = 50%

Safety Adjustments:
  + Can buy & sell: +15%
  + Low honeypot score (<0.1): +10%
  + Liquidity locked: +8%
  + Owner renounced: +7%
  + No restrictions: +5%

Volume Adjustments:
  + Good DEX volume: +10%
  + Social activity: +8%

Momentum Adjustments:
  + Rising velocity: +7%
  - Falling velocity: -10%

Final Win Probability = clamp(sum, 0, 1)

Action:
  if winProb >= 0.80 and autoExecute: BUY
  if winProb >= 0.80: LIST
  if winProb >= 0.60: MONITOR
  else: SKIP
```

### 6. CandidateListingAgent
**Purpose**: Manage queue of trading candidates

**Responsibilities**:
- Store candidate tokens with all reports
- Maintain status (pending/approved/rejected/executed)
- Provide queue for execution
- Track candidate history

**Output**: Candidate queue channel

### 7. ExecutionAgent
**Purpose**: Execute trades on blockchain

**Responsibilities**:
- Dry-run simulation before execution
- Build and sign transactions
- Broadcast to network
- Wait for confirmations
- Handle errors and retries
- Support multiple signing methods (OKX Wallet SDK, private key)

**Output**: `ExecutionResult` with transaction hash and status

**Chain-Specific**:
- **Base (EVM)**: Use router contracts, manage gas/nonce
- **Solana**: Use swap programs, manage token accounts

### 8. RiskManagerAgent
**Purpose**: Enforce risk controls and circuit breakers

**Responsibilities**:
- Check single position limit (1% default)
- Check total exposure limit (5% default)
- Track daily loss
- Trigger circuit breaker if limits exceeded
- Provide manual override capability
- Reset daily counters

**Output**: Approve/reject decisions

**Limits**:
- Single position: `amount <= balance * singlePositionPct`
- Total exposure: `total <= balance * totalExposurePct`
- Daily loss: `loss <= dailyLossLimit`

### 9. TelemetryAgent
**Purpose**: Track metrics and performance

**Responsibilities**:
- Record all agent activities
- Calculate performance metrics
- Provide metrics via API
- Periodic logging
- Support Prometheus export (future)

**Metrics**:
- Scanning: tokens scanned/found
- Filtering: tokens filtered/dropped
- Safety: checks performed, honeypots detected
- Strategy: evaluations, candidates listed
- Execution: trades executed, success/failure
- Financial: invested, profit, loss
- Performance: latency, execution time

## Data Flow

```
┌─────────────────┐
│ Blockchain      │
│ (Solana/Base)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ ChainScanner    │ ───► TokenFound
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ PreFilter       │ ───► PreFilteredToken
└────────┬────────┘
         │
         ├──────────────┐
         ▼              ▼
┌─────────────┐  ┌──────────────┐
│ OnChain     │  │ OffChain     │
│ Safety      │  │ Data         │
└──────┬──────┘  └──────┬───────┘
       │                │
       │    SafetyReport│
       │                │
       └────────┬───────┘
                │ OffChainMetrics
                ▼
       ┌─────────────────┐
       │ Strategy        │
       │ Evaluator       │
       └────────┬────────┘
                │ StrategyDecision
                ▼
       ┌─────────────────┐
       │ Candidate       │
       │ Listing         │
       └────────┬────────┘
                │
                ▼
       ┌─────────────────┐
       │ Risk            │ ◄── Approve/Reject
       │ Manager         │
       └────────┬────────┘
                │ Approved
                ▼
       ┌─────────────────┐
       │ Execution       │ ───► ExecutionResult
       └────────┬────────┘
                │
                ▼
       ┌─────────────────┐
       │ Telemetry       │ ───► Metrics/Logs
       └─────────────────┘
```

## Orchestrator

The `Orchestrator` coordinates all agents:

1. Initializes all agents
2. Starts ChainScannerAgent
3. Processes tokens through pipeline
4. Manages execution queue
5. Handles graceful shutdown

**Pipeline Processing**:
```go
token := <-scanner.TokenChannel
  ↓
filtered := prefilter.Filter(token)
  ↓
safety := safety.Evaluate(filtered)
  ↓
offchain := offchain.Gather(filtered)
  ↓
decision := strategy.Evaluate(safety, offchain)
  ↓
if decision.Action == "list" || "buy":
    candidate := listing.AddCandidate(...)
    ↓
    if autoExecute && risk.CanExecute():
        result := execution.Execute(candidate)
```

## Configuration

All agents use a shared `Config` object:
- Loaded from environment variables
- `.env` file support
- Type-safe with defaults
- Validation on load

## Error Handling

- Each agent handles its own errors
- Errors logged but don't crash system
- Failed tokens dropped from pipeline
- Telemetry tracks failure rates
- Circuit breaker for excessive failures

## Concurrency

- Scanner runs on separate goroutines per chain
- Token processing is concurrent (goroutine per token)
- Execution queue processed serially to avoid conflicts
- All shared state protected by mutexes
- Channels for inter-agent communication

## Testing Strategy

1. **Unit Tests**: Test individual agent logic
2. **Integration Tests**: Test agent interactions
3. **Dry-Run Testing**: Test full system without execution
4. **Small Amount Testing**: Test with minimal funds
5. **Monitoring**: Continuous validation in production

## Scalability

**Current Design**:
- Single process
- Suitable for monitoring 2 chains
- ~1000 tokens/hour capacity

**Future Scaling Options**:
- Distributed agents via message queue
- Separate scanner services per chain
- Database for candidate persistence
- Redis for shared state
- Kubernetes deployment

## Security

1. **Private Key Protection**
   - Never in code
   - Environment variables or KMS
   - OKX Wallet SDK preferred

2. **RPC Rate Limiting**
   - Respect endpoint limits
   - Exponential backoff
   - Multiple provider fallback

3. **Input Validation**
   - Validate all on-chain data
   - Sanitize addresses
   - Check ranges

4. **Audit Trail**
   - Log all decisions
   - Record all trades
   - Preserve evidence

## Monitoring

**Key Metrics**:
- Tokens/second scanned
- Candidates/hour listed
- Win rate (profitable trades %)
- Average ROI
- Circuit breaker triggers
- API latency

**Alerts**:
- Circuit breaker triggered
- Execution failures
- RPC connection issues
- Unusual activity patterns

## Future Enhancements

1. **Machine Learning**
   - Train models on historical data
   - Improve win probability predictions
   - Pattern recognition

2. **Multi-Chain Support**
   - Ethereum, Polygon, Arbitrum
   - Cross-chain arbitrage

3. **Advanced Strategies**
   - Market making
   - Liquidity provision
   - Arbitrage opportunities

4. **Social Integration**
   - Discord/Telegram notifications
   - Community signals
   - Influencer tracking

5. **Portfolio Management**
   - Automated rebalancing
   - Position tracking
   - Tax reporting
