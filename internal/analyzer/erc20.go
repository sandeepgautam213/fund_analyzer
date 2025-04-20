package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan" // Internal package to interact with Etherscan API
)

// ParseERC20Transfers fetches and parses all ERC-20 token transfers for a given address.
// It returns a list of Beneficiaries, each with a list of transfers and total amount sent.

func ParseERC20Transfers(address string) []Beneficiary {
	// Fetch raw JSON data of ERC-20 transfers from Etherscan API
	tokenRaw, err := etherscan.FetchTxs(address, "tokentx")
	if err != nil {
		log.Println("Error fetching ERC-20 transfers:", err)
		return nil
	}
	// Define a variable to hold the parsed JSON data
	var tokenData TokenTxResult
	// Unmarshal (parse) the JSON response into tokenData struct
	if err := json.Unmarshal(tokenRaw, &tokenData); err != nil {
		log.Println("Error parsing ERC-20 transfers:", err)
		return nil // Return nil if JSON parsing fails
	}

	// beneficiariesMap groups transfers by "recipient + token"

	beneficiariesMap := make(map[string][]TxInfo)
	// totals stores the total amount sent to each recipient per token
	totals := make(map[string]float64)
	// Loop over each token transfer in the result
	for _, tx := range tokenData.Result {
		// Skip if transfer has no recipient or it's a self-transfer
		if tx.To == "" || tx.To == address {
			continue
		}
		// Convert transfer amount from string to float64

		tokenAmount, _ := strconv.ParseFloat(tx.Value, 64)
		// Convert token decimal value to integer
		decimals, _ := strconv.Atoi(tx.TokenDecimal)
		// Normalize token amount using token decimals (i.e., divide by 10^decimals)
		tokenAmount /= pow10(decimals)
		// Convert timestamp (string) to Unix time (int64)
		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		// Format timestamp into human-readable format
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
		// Create a transaction info object
		txinfo := TxInfo{
			TxAmount:      tokenAmount,
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}
		// Use recipient address + token symbol as a key to group by token & address
		key := fmt.Sprintf("%s (token: %s)", tx.To, tx.TokenSymbol)
		// Append the transaction to the recipient's list
		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		// Add to the total amount sent to this recipient for this token
		totals[key] += tokenAmount
	}
	// Final output list of Beneficiary structs
	var result []Beneficiary
	// Build each Beneficiary from the maps
	for addr, txs := range beneficiariesMap {
		result = append(result, Beneficiary{
			BeneficiaryAddress: addr,
			Amount:             totals[addr],
			Transactions:       txs,
		})
	}
	// Return the full result list
	return result
}

// pow10 returns 10^n as float64
func pow10(n int) float64 {
	res := 1.0
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}
