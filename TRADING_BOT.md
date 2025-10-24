# Meme Coin Trading Bot

A sophisticated automated trading system for scanning and trading newly issued meme coins on Solana and Base chains with comprehensive risk management and honeypot detection.

## Overview

This trading bot implements a multi-agent architecture that:
- Continuously monitors Solana and Base chains for new token creation
- Performs comprehensive safety checks (honeypot detection, buy/sell simulation)
- Gathers off-chain metrics (DEX/CEX volume, social media signals)
- Evaluates trading opportunities using strategy algorithms
- Executes trades when win probability ≥ 80% (configurable)
- Implements robust risk management and circuit breakers

## Architecture

### Agent System

The bot consists of 9 specialized agents working together:

1. **ChainScannerAgent** - Monitors on-chain events for new token creation
2. **PreFilterAgent** - Applies basic filtering rules (blacklist, liquidity checks)
3. **OnChainSafetyAgent** - Performs honeypot detection and buy/sell simulation
4. **OffChainDataAgent** - Gathers trading volume and social metrics
5. **StrategyEvaluatorAgent** - Calculates win probability and recommends actions
6. **CandidateListingAgent** - Manages queue of trading candidates
7. **ExecutionAgent** - Executes trades via OKX Wallet SDK or private key
8. **RiskManagerAgent** - Enforces position limits and circuit breakers
9. **TelemetryAgent** - Tracks metrics and performance

### Data Flow

```
New Token Created
    ↓
ChainScannerAgent (discovers token)
    ↓
PreFilterAgent (basic filtering)
    ↓
OnChainSafetyAgent (honeypot check)
    ↓
OffChainDataAgent (volume & social metrics)
    ↓
StrategyEvaluatorAgent (win probability ≥ 80%)
    ↓
CandidateListingAgent (queue for execution)
    ↓
RiskManagerAgent (check limits)
    ↓
ExecutionAgent (execute trade)
    ↓
TelemetryAgent (record metrics)
```

## Quick Start

### Prerequisites

- Go 1.20 or higher
- Access to Solana and Base RPC nodes
- OKX Wallet SDK or private key (for live trading)
- API keys for CoinGecko, Twitter (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot
```

2. Configure environment:
```bash
cp .env.example .env
# Edit .env with your settings
```

3. Build the trading bot:
```bash
make build-backend
```

### Running in Dry-Run Mode (Recommended First)

```bash
# Edit .env and ensure:
# DRY_RUN=true
# AUTO_EXECUTE=false

make run-trading
```

This will:
- Monitor chains for new tokens
- Evaluate opportunities
- Log actions without executing real trades
- Build confidence in the system

### Running with Auto-Execution

⚠️ **WARNING**: Only use with small amounts you can afford to lose!

```bash
# Edit .env:
# DRY_RUN=false
# AUTO_EXECUTE=true
# ACCOUNT_BALANCE=1000  # Your actual balance

make run-trading
```

## Configuration

### Essential Settings

```bash
# General
DRY_RUN=true                    # Test mode (no real trades)
AUTO_EXECUTE=false              # Auto-execute trades

# Strategy
WIN_PROBABILITY_THRESHOLD=0.80  # Minimum 80% win probability
MIN_VOLUME_DEX=10000.0          # Minimum DEX volume
MIN_LIQUIDITY=5000.0            # Minimum initial liquidity
MAX_HONEYPOT_SCORE=0.2          # Maximum honeypot risk (0-1)

# Risk Management
SINGLE_POSITION_PCT=0.01        # Max 1% per trade
TOTAL_EXPOSURE_PCT=0.05         # Max 5% total exposure
DAILY_LOSS_LIMIT=500.0          # Daily loss limit (USD)
ACCOUNT_BALANCE=10000.0         # Your account balance
```

### RPC Endpoints

```bash
# Solana
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com

# Base
BASE_RPC_URL=https://mainnet.base.org
BASE_WS_URL=wss://mainnet.base.org
```

### API Keys (Optional but Recommended)

```bash
COINGECKO_API_KEY=your_key_here
TWITTER_API_KEY=your_key_here
OKX_API_KEY=your_key_here
```

### Wallet Configuration

```bash
# Option 1: OKX Wallet SDK (Recommended)
USE_OKX_WALLET=true

# Option 2: Private Key (Use with caution!)
USE_OKX_WALLET=false
PRIVATE_KEY=your_private_key_here  # Never commit this!
```

## API Endpoints

The trading bot exposes a REST API on port 8080:

### Health Check
```bash
GET /api/health
```

### System Status
```bash
GET /api/status
Response: {
  "status": "running",
  "candidate_count": 5,
  "trading_halted": false,
  "metrics": {...}
}
```

### View Candidates
```bash
GET /api/candidates
Response: {
  "count": 5,
  "candidates": [...]
}
```

### Metrics
```bash
GET /api/metrics
Response: {
  "tokens_scanned": 150,
  "tokens_found": 45,
  "candidates_listed": 12,
  "trades_executed": 3,
  "total_profit": 250.50,
  ...
}
```

### Risk Status
```bash
GET /api/risk
Response: {
  "current_exposure": 450.00,
  "daily_loss": 125.50,
  "trading_halted": false,
  ...
}
```

### Resume Trading (Manual Override)
```bash
POST /api/risk/resume
```

## Safety Features

### Honeypot Detection

The bot performs comprehensive safety checks:

1. **Simulated Buy/Sell** - Tests if token can actually be sold
2. **Transfer Restrictions** - Checks for blacklist/whitelist mechanisms
3. **Owner Controls** - Verifies if owner is renounced
4. **Tax Analysis** - Detects excessive transaction taxes
5. **Liquidity Lock** - Confirms liquidity is locked
6. **Slippage Check** - Ensures slippage is within acceptable range

A token passes safety if:
- `can_buy == true && can_sell == true`
- `honeypot_score < 0.2` (configurable)
- Simulated sell succeeds with acceptable slippage

### Risk Management

1. **Position Limits**
   - Single position: max 1% of balance (default)
   - Total exposure: max 5% of balance (default)

2. **Circuit Breaker**
   - Automatically halts trading if daily loss limit reached
   - Requires manual resume via API

3. **Daily Reset**
   - Loss counters reset every 24 hours
   - Exposure tracking maintained

## Strategy Evaluation

### Win Probability Calculation

The bot calculates win probability based on:

**Safety Factors (Most Important)**
- Can buy and sell: +15%
- Low honeypot score (<0.1): +10%
- Liquidity locked: +8%
- Owner renounced: +7%
- No transfer restrictions: +5%

**Volume Factors**
- Good DEX volume: +10%
- Social activity: +8%

**Momentum Factors**
- Rising velocity: +7%
- Falling velocity: -10%

**Base probability**: 50%
**Threshold**: 80% (configurable)

### Action Determination

- **buy** - WinProb ≥ 80% and AUTO_EXECUTE=true
- **list** - WinProb ≥ 80% and AUTO_EXECUTE=false
- **monitor** - WinProb ≥ 60%
- **skip** - WinProb < 60%

## Monitoring

### Metrics

The TelemetryAgent tracks:
- Tokens scanned/found/filtered
- Safety checks performed
- Honeypots detected
- Candidates listed
- Trades executed (success/failed)
- Financial performance (invested, profit, loss)
- Performance (decision latency, execution time)

### Logging

Logs are output to stdout with levels:
- INFO - Normal operations
- WARNING - Potential issues
- ERROR - Failures

## Chain-Specific Implementation

### Base (EVM)

- Uses `eth_call` for simulation
- Monitors Uniswap factory events
- Checks ERC20 contract patterns
- Gas estimation and management
- Nonce tracking

### Solana

- Uses `simulateTransaction` RPC
- Monitors SPL Token programs
- Checks Raydium/Orca pools
- Token account management
- Rent-exempt requirements

## Testing

### Run Tests
```bash
make test
```

### Dry-Run Testing

Always test in dry-run mode first:
1. Set `DRY_RUN=true`
2. Monitor for several hours
3. Review candidate selections
4. Verify strategy logic
5. Check risk controls

### Small Amount Testing

When ready for live trading:
1. Set `ACCOUNT_BALANCE=100` (small amount)
2. Set `SINGLE_POSITION_PCT=0.05` (5%)
3. Set `DAILY_LOSS_LIMIT=10`
4. Monitor closely for first 24 hours

## Security Best Practices

1. **Never Commit Private Keys**
   - Use environment variables
   - Consider KMS/Vault for production

2. **Use OKX Wallet SDK**
   - Preferred over raw private keys
   - Better security model

3. **Start with Testnet**
   - Test on Solana devnet
   - Test on Base testnet

4. **Gradual Scaling**
   - Start with minimal amounts
   - Increase gradually after proven track record

5. **Monitor Continuously**
   - Set up alerting
   - Review daily metrics
   - Check for anomalies

## Compliance & Legal

⚠️ **IMPORTANT DISCLAIMERS**

1. **High Risk**: Cryptocurrency trading carries substantial risk of loss
2. **No Guarantees**: Past performance doesn't guarantee future results
3. **Meme Coins**: Extremely volatile and speculative
4. **Regulatory**: Ensure compliance with local laws
5. **Tax Implications**: Track all trades for tax reporting
6. **Use at Your Own Risk**: This is experimental software

## Troubleshooting

### Trading Not Executing

Check:
1. `AUTO_EXECUTE=true` in .env
2. Risk limits not exceeded
3. Trading not halted (check `/api/risk`)
4. RPC endpoints accessible
5. Wallet configured correctly

### No Tokens Found

Check:
1. RPC endpoints working
2. Scan intervals configured
3. Network connectivity
4. Check logs for errors

### High Rejection Rate

Adjust:
1. Lower `WIN_PROBABILITY_THRESHOLD`
2. Lower `MIN_LIQUIDITY`
3. Lower `MIN_VOLUME_DEX`
4. Check blacklist settings

## Development

### Adding New Chains

1. Add chain constant to `pkg/models/events.go`
2. Implement scanner in `pkg/agents/scanner/`
3. Implement safety checks in `pkg/agents/safety/`
4. Implement execution in `pkg/agents/execution/`
5. Add configuration in `pkg/config/`

### Extending Strategy

Modify `pkg/agents/strategy/strategy.go`:
- Add new factors to win probability
- Adjust confidence calculation
- Customize position sizing
- Add ML models

### Custom Agents

Create new agent in `pkg/agents/`:
1. Define agent struct
2. Implement evaluation logic
3. Integrate with orchestrator

## Support

For issues or questions:
- Open an issue on GitHub
- Review documentation
- Check API status endpoints

## License

MIT License - See LICENSE file

---

**Remember**: Always start in dry-run mode and only trade with amounts you can afford to lose!
