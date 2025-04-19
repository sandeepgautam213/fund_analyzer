package analyzer

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(2)

	payerMap := make(map[string][]TxInfo)
	payerTotals := make(map[string]float64)
	var mu sync.Mutex

	go func() {
		defer wg.Done()
		normal := parsePayerTxs(address, "txlist")
		mu.Lock()
		mergePayers(payerMap, payerTotals, normal)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		internal := parsePayerTxs(address, "txlistinternal")
		mu.Lock()
		mergePayers(payerMap, payerTotals, internal)
		mu.Unlock()
	}()

	wg.Wait()

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

func parsePayerTxs(address string, action string) []Payer {
	raw, err := etherscan.FetchTxs(address, action)
	if err != nil {
		log.Println("Error fetching ", action, "transactions:", err)
		return nil
	}

	var data TxResult
	json.Unmarshal(raw, &data)

	payerMap := make(map[string][]TxInfo)
	payerTotals := make(map[string]float64)

	for _, tx := range data.Result {
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

	var result []Payer
	for from, txs := range payerMap {
		result = append(result, Payer{
			PayerAddress: from,
			Amount:       payerTotals[from],
			Transactions: txs,
		})
	}

	return result
}

func mergePayers(dest map[string][]TxInfo, totals map[string]float64, incoming []Payer) {
	for _, p := range incoming {
		dest[p.PayerAddress] = append(dest[p.PayerAddress], p.Transactions...)
		totals[p.PayerAddress] += p.Amount
	}
}
