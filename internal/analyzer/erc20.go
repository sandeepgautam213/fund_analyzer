package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

func ParseERC20Transfers(address string) []Beneficiary {
	tokenRaw, err := etherscan.FetchTxs(address, "tokentx")
	if err != nil {
		log.Println("Error fetching ERC-20 transfers:", err)
		return nil
	}

	var tokenData TokenTxResult
	if err := json.Unmarshal(tokenRaw, &tokenData); err != nil {
		log.Println("Error parsing ERC-20 transfers:", err)
		return nil
	}

	beneficiariesMap := make(map[string][]TxInfo)
	totals := make(map[string]float64)

	for _, tx := range tokenData.Result {
		if tx.To == "" || tx.To == address {
			continue
		}

		tokenAmount, _ := strconv.ParseFloat(tx.Value, 64)
		decimals, _ := strconv.Atoi(tx.TokenDecimal)
		tokenAmount /= pow10(decimals)

		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      tokenAmount,
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}

		key := fmt.Sprintf("%s (token: %s)", tx.To, tx.TokenSymbol)
		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		totals[key] += tokenAmount
	}

	var result []Beneficiary
	for addr, txs := range beneficiariesMap {
		result = append(result, Beneficiary{
			BeneficiaryAddress: addr,
			Amount:             totals[addr],
			Transactions:       txs,
		})
	}

	return result
}

func pow10(n int) float64 {
	res := 1.0
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}
