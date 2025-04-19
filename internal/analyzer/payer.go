package analyzer

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

type Payer struct {
	PayerAddress string   `json:"beneficiary_address"`
	Amount       float64  `json:"amount"`
	Date         string   `json:"date,omitempty"`
	Transactions []TxInfo `json:"transactions"`
}

func AnalyzePayers(address string) []Payer {
	normalRaw, _ := etherscan.FetchTxs(address, "txlist")
	internalRaw, _ := etherscan.FetchTxs(address, "txlistinternal")

	var normalData, internalData TxResult
	json.Unmarshal(normalRaw, &normalData)
	json.Unmarshal(internalRaw, &internalData)

	allTxs := append(normalData.Result, internalData.Result...)

	payerMap := make(map[string][]TxInfo)
	payerTotals := make(map[string]float64)

	for _, tx := range allTxs {
		if strings.ToLower(tx.To) != strings.ToLower(address) {
			continue
		}

		amt, _ := strconv.ParseFloat(tx.Value, 64)
		amt /= 1e18

		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		formatted := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      amt,
			DateTime:      formatted,
			TransactionID: tx.Hash,
		}

		payerMap[tx.From] = append(payerMap[tx.From], txinfo)
		payerTotals[tx.From] += amt
	}

	var payers []Payer
	for from, txs := range payerMap {

		earliest := txs[0].DateTime
		for _, tx := range txs {
			if tx.DateTime < earliest {
				earliest = tx.DateTime
			}
		}

		payers = append(payers, Payer{
			PayerAddress: from,
			Amount:       payerTotals[from],
			Date:         earliest,
			Transactions: txs,
		})
	}

	return payers
}
