package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

type Transaction struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	TimeStamp string `json:"timeStamp"`
}

type TxResult struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type TxInfo struct {
	TxAmount      float64 `json:"tx_amount"`
	DateTime      string  `json:"date_time"`
	TransactionID string  `json:"transaction_id"`
}

type Beneficiary struct {
	BeneficiaryAddress string   `json:"beneficiary_address"`
	Amount             float64  `json:"amount"`
	Transactions       []TxInfo `json:"transactions"`
}
type TokenTransfer struct {
	Hash         string `json:"hash"`
	From         string `json:"from"`
	To           string `json:"to"`
	Value        string `json:"value"`
	TimeStamp    string `json:"timeStamp"`
	TokenName    string `json:"tokenName"`
	TokenSymbol  string `json:"tokenSymbol"`
	TokenDecimal string `json:"tokenDecimal"`
}
type TokenTxResult struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  []TokenTransfer `json:"result"`
}

type NFTTransfer struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	TokenID     string `json:"tokenID"`
	TimeStamp   string `json:"timeStamp"`
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
}

type NFTTxResult struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []NFTTransfer `json:"result"`
}

func AnalyzeAddress(address string) []Beneficiary {
	// Fetch normal transactions
	normalRaw, err := etherscan.FetchTxs(address, "txlist")
	if err != nil {
		log.Println("Error fetching normal transactions:", err)
		return nil
	}

	// Fetch internal transactions
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

	// Fetch ERC-20 token transfers
	tokenRaw, err := etherscan.FetchTxs(address, "tokentx")
	if err != nil {
		log.Println("Error fetching token transfers:", err)
	}

	var tokenData TokenTxResult
	if err := json.Unmarshal(tokenRaw, &tokenData); err != nil {
		log.Println("Error parsing token transfers:", err)
	}

	// Fetch ERC-721 NFT transfers
	nftRaw, err := etherscan.FetchTxs(address, "tokennfttx")
	if err != nil {
		log.Println("Error fetching NFT transfers:", err)
	}

	var nftData NFTTxResult
	if err := json.Unmarshal(nftRaw, &nftData); err != nil {
		log.Println("Error parsing NFT transfers:", err)
	}

	// Merge all transactions
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

	// Process token transfers
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

		// Use a combined key to differentiate token vs ETH
		key := fmt.Sprintf("%s (token: %s)", tx.To, tx.TokenSymbol)

		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		totals[key] += tokenAmount
	}

	for _, tx := range nftData.Result {
		if tx.To == "" || tx.To == address {
			continue
		}

		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      1, // each NFT is unique
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}

		key := fmt.Sprintf("%s (NFT: %s #%s)", tx.To, tx.TokenSymbol, tx.TokenID)

		// log.Printf("NFT ➜ To: %s, Token: %s #%s, Tx: %s", tx.To, tx.TokenSymbol, tx.TokenID, tx.Hash)

		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		totals[key] += 1 // each transfer is 1 NFT
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
