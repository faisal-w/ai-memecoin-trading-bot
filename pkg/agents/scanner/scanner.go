package scanner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mumugogoing/meme_bot/pkg/config"
	"github.com/mumugogoing/meme_bot/pkg/models"
)

// ChainScannerAgent monitors on-chain events for new tokens
type ChainScannerAgent struct {
	config       *config.Config
	tokenChannel chan models.TokenFound
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewChainScannerAgent creates a new chain scanner agent
func NewChainScannerAgent(cfg *config.Config) *ChainScannerAgent {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChainScannerAgent{
		config:       cfg,
		tokenChannel: make(chan models.TokenFound, 100),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins scanning both chains
func (s *ChainScannerAgent) Start() {
	log.Println("ChainScannerAgent: Starting chain monitoring...")
	
	// Start Solana scanner
	s.wg.Add(1)
	go s.scanSolana()
	
	// Start Base scanner
	s.wg.Add(1)
	go s.scanBase()
	
	log.Println("ChainScannerAgent: Chain monitoring started")
}

// Stop stops all scanning operations
func (s *ChainScannerAgent) Stop() {
	log.Println("ChainScannerAgent: Stopping...")
	s.cancel()
	s.wg.Wait()
	close(s.tokenChannel)
	log.Println("ChainScannerAgent: Stopped")
}

// GetTokenChannel returns the channel for discovered tokens
func (s *ChainScannerAgent) GetTokenChannel() <-chan models.TokenFound {
	return s.tokenChannel
}

// scanSolana monitors Solana chain for new tokens
func (s *ChainScannerAgent) scanSolana() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.config.ScanIntervalSolana)
	defer ticker.Stop()
	
	log.Printf("ChainScannerAgent: Solana scanner started (interval: %v)\n", s.config.ScanIntervalSolana)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.scanSolanaNewTokens()
		}
	}
}

// scanBase monitors Base chain for new tokens
func (s *ChainScannerAgent) scanBase() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.config.ScanIntervalBase)
	defer ticker.Stop()
	
	log.Printf("ChainScannerAgent: Base scanner started (interval: %v)\n", s.config.ScanIntervalBase)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.scanBaseNewTokens()
		}
	}
}

// Solana RPC Response types
type solanaRPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

type solanaBlockData struct {
	Blockhash         string `json:"blockhash"`
	PreviousBlockhash string `json:"previousBlockhash"`
	ParentSlot        int64  `json:"parentSlot"`
	Transactions      []struct {
		Transaction struct {
			Message struct {
				Instructions []struct {
					ProgramIDIndex int           `json:"programIdIndex"`
					Accounts       []int         `json:"accounts"`
					Data           string        `json:"data"`
				} `json:"instructions"`
			} `json:"message"`
		} `json:"transaction"`
		Meta struct {
			InnerInstructions []struct {
				Index        int `json:"index"`
				Instructions []struct {
					ProgramIDIndex int    `json:"programIdIndex"`
					Accounts       []int  `json:"accounts"`
					Data           string `json:"data"`
				} `json:"instructions"`
			} `json:"innerInstructions"`
			LogMessages []string `json:"logMessages"`
		} `json:"meta"`
	} `json:"transactions"`
}

// Base/EVM RPC Response types
type evmBlockData struct {
	Hash       string        `json:"hash"`
	Number     string        `json:"number"`
	Timestamp  string        `json:"timestamp"`
	Difficulty string        `json:"difficulty"`
}

type evmLogEntry struct {
	Address  string   `json:"address"`
	Topics   []string `json:"topics"`
	Data     string   `json:"data"`
	BlockNum string   `json:"blockNumber"`
	TxHash   string   `json:"transactionHash"`
}

// scanSolanaNewTokens scans for new tokens on Solana by fetching the latest block
func (s *ChainScannerAgent) scanSolanaNewTokens() {
	log.Println("ChainScannerAgent: Scanning Solana for new tokens...")

	// Get latest slot
	slot, err := s.getSolanaSlot()
	if err != nil {
		log.Printf("ChainScannerAgent: Error getting Solana slot: %v\n", err)
		return
	}

	// Get block data with full transaction details
	blockData, err := s.getSolanaBlockData(slot)
	if err != nil {
		log.Printf("ChainScannerAgent: Error getting Solana block data: %v\n", err)
		return
	}

	// Process transactions to find token creation and pool events
	s.processSolanaTransactions(blockData)
}

// getSolanaSlot fetches the latest confirmed slot
func (s *ChainScannerAgent) getSolanaSlot() (int64, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSlot",
		"params":  []interface{}{"confirmed"},
	}

	var slot int64
	if err := s.makeJSONRPCCall(s.config.SolanaRPCURL, reqBody, &slot); err != nil {
		return 0, err
	}
	return slot, nil
}

// getSolanaBlockData fetches full block data including all transactions
func (s *ChainScannerAgent) getSolanaBlockData(slot int64) (*solanaBlockData, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlock",
		"params": []interface{}{
			slot,
			map[string]interface{}{
				"encoding":                       "json",
				"transactionDetails":             "full",
				"rewards":                        false,
				"maxSupportedTransactionVersion": 0,
			},
		},
	}

	var blockData solanaBlockData
	if err := s.makeJSONRPCCall(s.config.SolanaRPCURL, reqBody, &blockData); err != nil {
		return nil, err
	}
	return &blockData, nil
}

// processSolanaTransactions extracts token creation and liquidity pool events
func (s *ChainScannerAgent) processSolanaTransactions(blockData *solanaBlockData) {
	// Known program IDs
	const (
		splTokenProgramID = "TokenkegQfeZyiNwAJsyFbPVwwQQfzZhy7zJES69Gm9"
		raydiumProgram    = "675kPX9MHTjS2zt1qrXjVVxt2Y8F4We6g76bfeKFVgJ"
		orcaProgram       = "9W959DqBcxTVd6bPV9vA75Pwbnx2zktxQJ5pWKeS6L8M"
	)

	for _, tx := range blockData.Transactions {
		// Check for SPL Token InitializeMint instructions (token creation)
		for _, instr := range tx.Transaction.Message.Instructions {
			// Check program index for token creation
			if instr.ProgramIDIndex < len(tx.Transaction.Message.Instructions) {
				// Parse instruction data for token creation
				// InitializeMint is instruction 0 in SPL Token program
				if strings.Contains(fmt.Sprintf("%d", instr.ProgramIDIndex), "0") {
					s.parseTokenCreationInstruction(&tx, instr.Data)
				}
			}
		}

		// Check inner instructions for liquidity pool creation events
		if tx.Meta.InnerInstructions != nil {
			for _, innerInstr := range tx.Meta.InnerInstructions {
				for _, instr := range innerInstr.Instructions {
					// Look for Raydium/Orca pool creation
					if instr.ProgramIDIndex == 0 || instr.ProgramIDIndex == 1 {
						s.parseLiquidityPoolCreation(&tx, instr.Data)
					}
				}
			}
		}

		// Parse log messages for events
		if tx.Meta.LogMessages != nil {
			s.parseLogsForEvents(tx.Meta.LogMessages)
		}
	}
}

// parseTokenCreationInstruction extracts token info from InitializeMint instructions
func (s *ChainScannerAgent) parseTokenCreationInstruction(tx interface{}, instructionData string) {
	// In production, decode base58-encoded instruction data and extract:
	// - Token mint address
	// - Decimals
	// - Owner/creator
	log.Printf("ChainScannerAgent: Found potential token creation instruction\n")
}

// parseLiquidityPoolCreation extracts liquidity pool creation events
func (s *ChainScannerAgent) parseLiquidityPoolCreation(tx interface{}, instructionData string) {
	// In production, decode instruction data to find:
	// - Pool address
	// - Token pair
	// - Initial reserves
	log.Printf("ChainScannerAgent: Found potential liquidity pool creation\n")
}

// parseLogsForEvents extracts events from transaction logs
func (s *ChainScannerAgent) parseLogsForEvents(logs []string) {
	for _, log := range logs {
		// Look for token creation and pool creation events
		if strings.Contains(log, "InitializeMint") {
			log.Printf("ChainScannerAgent: Token initialization detected in logs\n")
		}
		if strings.Contains(log, "AddLiquidity") || strings.Contains(log, "PoolCreated") {
			log.Printf("ChainScannerAgent: Liquidity event detected in logs\n")
		}
	}
}

// scanBaseNewTokens scans for new tokens on Base (EVM) by fetching the latest block
func (s *ChainScannerAgent) scanBaseNewTokens() {
	log.Println("ChainScannerAgent: Scanning Base for new tokens...")

	// Get latest block number
	blockNum, err := s.getBaseBlockNumber()
	if err != nil {
		log.Printf("ChainScannerAgent: Error getting Base block number: %v\n", err)
		return
	}

	// Get logs for Uniswap V2 factory PairCreated events
	logs, err := s.getBaseLogs(blockNum)
	if err != nil {
		log.Printf("ChainScannerAgent: Error getting Base logs: %v\n", err)
		return
	}

	// Process logs to extract token pairs
	s.processBaseLogs(logs)
}

// getBaseBlockNumber fetches the latest block number
func (s *ChainScannerAgent) getBaseBlockNumber() (string, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_blockNumber",
	}

	var blockNum string
	if err := s.makeJSONRPCCall(s.config.BaseRPCURL, reqBody, &blockNum); err != nil {
		return "", err
	}
	return blockNum, nil
}

// getBaseLogs fetches logs for the latest block looking for Uniswap factory events
func (s *ChainScannerAgent) getBaseLogs(blockNum string) ([]evmLogEntry, error) {
	// Uniswap V2 Factory address on Base: 0x8909Dc15e40953b386fa8f440dffe6e6543f20df
	// PairCreated event signature: keccak256("PairCreated(address,address,address,uint256)")
	const (
		uniswapFactoryBase = "0x8909Dc15e40953b386fa8f440dffe6e6543f20df"
		pairCreatedTopic   = "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"
	)

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getLogs",
		"params": []interface{}{
			map[string]interface{}{
				"fromBlock": blockNum,
				"toBlock":   blockNum,
				"address":   uniswapFactoryBase,
				"topics":    []interface{}{pairCreatedTopic},
			},
		},
	}

	var logs []evmLogEntry
	if err := s.makeJSONRPCCall(s.config.BaseRPCURL, reqBody, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// processBaseLogs extracts token pairs from Uniswap factory logs
func (s *ChainScannerAgent) processBaseLogs(logs []evmLogEntry) {
	for _, logEntry := range logs {
		// PairCreated event has indexed parameters: token0, token1, pair, and index
		// Topics[1] = token0, Topics[2] = token1, Topics[3] = pair
		if len(logEntry.Topics) >= 4 {
			token0 := "0x" + logEntry.Topics[1][len(logEntry.Topics[1])-40:]
			token1 := "0x" + logEntry.Topics[2][len(logEntry.Topics[2])-40:]
			pairAddr := "0x" + logEntry.Topics[3][len(logEntry.Topics[3])-40:]

			log.Printf("ChainScannerAgent: Uniswap pair created - %s + %s = %s\n", token0, token1, pairAddr)

			// Emit token found events for both tokens
			s.emitTokenFound(models.TokenFound{
				Chain:          models.ChainBase,
				TokenAddress:   token0,
				FirstSeenTS:    time.Now().Unix(),
				CreatorAddress: pairAddr,
				TxHash:         logEntry.TxHash,
				Metadata: map[string]string{
					"pair_address": pairAddr,
					"token_type":   "token0",
				},
			})

			s.emitTokenFound(models.TokenFound{
				Chain:          models.ChainBase,
				TokenAddress:   token1,
				FirstSeenTS:    time.Now().Unix(),
				CreatorAddress: pairAddr,
				TxHash:         logEntry.TxHash,
				Metadata: map[string]string{
					"pair_address": pairAddr,
					"token_type":   "token1",
				},
			})
		}
	}
}

// makeJSONRPCCall makes a JSON-RPC call to a blockchain RPC endpoint
func (s *ChainScannerAgent) makeJSONRPCCall(rpcURL string, reqBody map[string]interface{}, result interface{}) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// Marshal request body
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make RPC call: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal JSON-RPC response
	var rpcResp solanaRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for RPC errors
	if rpcResp.Error != nil {
		return fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	// Unmarshal result
	if err := json.Unmarshal(rpcResp.Result, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return nil
}

// emitTokenFound sends a discovered token to the channel
func (s *ChainScannerAgent) emitTokenFound(token models.TokenFound) {
	select {
	case s.tokenChannel <- token:
		log.Printf("ChainScannerAgent: Token found - %s on %s\n", token.TokenAddress, token.Chain)
	case <-s.ctx.Done():
		return
	default:
		log.Println("ChainScannerAgent: Warning - token channel full, dropping event")
	}
}