package analyzer

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

func ParseETHTransfers(address string) []Beneficiary {
	normalRaw, err := etherscan.FetchTxs(address, "txlist")
	if err != nil {
		log.Println("Error fetching normal transactions:", err)
		return nil
	}

	internalRaw, err := etherscan.FetchTxs(address, "txlistinternal")
	if err != nil {
		log.Println("Error fetching internal transactions:", err)
		return nil
	}

	var normalData, internalData TxResult
	if err := json.Unmarshal(normalRaw, &normalData); err != nil {
		log.Println("Failed to parse normal txs:", err)
		return nil
	}
	if err := json.Unmarshal(internalRaw, &internalData); err != nil {
		log.Println("Failed to parse internal txs:", err)
		return nil
	}

	allTxs := append(normalData.Result, internalData.Result...)
	beneficiariesMap := make(map[string][]TxInfo)
	totals := make(map[string]float64)

	for _, tx := range allTxs {
		if tx.To == "" || tx.To == address {
			continue
		}

		ethAmount, _ := strconv.ParseFloat(tx.Value, 64)
		ethAmount /= 1e18

		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      ethAmount,
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}

		beneficiariesMap[tx.To] = append(beneficiariesMap[tx.To], txinfo)
		totals[tx.To] += ethAmount
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
