package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

type ERC1155Log struct {
	Address         string   `json:"address"`
	Topics          []string `json:"topics"`
	Data            string   `json:"data"`
	TimeStamp       string   `json:"timeStamp"`
	TransactionHash string   `json:"transactionHash"`
}

type LogResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Result  []ERC1155Log `json:"result"`
}

func FetchAndParseERC1155Logs(address string) []Beneficiary {
	var results []Beneficiary                     // Final list of beneficiaries
	beneficiariesMap := make(map[string][]TxInfo) // Map to group transactions by receiver and tokenID
	totals := make(map[string]float64)            // Map to store total tokens received by each group

	// OpenSea Shared Storefront - ERC-1155 contract
	erc1155Contract := "0x495f947276749ce646f68ac8c248420045cb7b5e"
	// Keccak-256 hash of TransferSingle event signature
	topicTransferSingle := "0xc3d58168c5bfaa58138de29b8634c2d14f7145e7f7a579f8b0b0b6b7f6fe15f6"

	// Fetch logs from Etherscan for the TransferSingle event
	logRaw, err := etherscan.FetchERC1155Logs(erc1155Contract, "0", "latest", topicTransferSingle)
	if err != nil {
		log.Println("Error fetching ERC-1155 logs:", err)
		return nil
	}
	// Unmarshal JSON into Go struct
	var logs LogResponse
	err = json.Unmarshal(logRaw, &logs)
	if err != nil {
		log.Println("Error unmarshaling ERC-1155 logs:", err)
		return nil
	}

	for _, l := range logs.Result {
		// Skip invalid logs with missing topics
		if len(l.Topics) < 4 {
			continue
		}

		// Extract to address from topic[3]
		toAddr := "0x" + l.Topics[2][26:]
		if strings.ToLower(toAddr) == strings.ToLower(address) {
			continue // skip self-transfer
		}
		// Convert the token ID from hex to big.Int
		tokenId := new(big.Int)
		tokenId.SetString(l.Topics[3][2:], 16)

		amount := new(big.Int)
		amount.SetString(l.Data[2:], 16)

		if amount.Sign() == 0 {
			continue
		}

		ts, _ := strconv.ParseInt(l.TimeStamp, 10, 64)
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      float64(amount.Int64()),
			DateTime:      formattedTime,
			TransactionID: l.TransactionHash,
		}

		key := fmt.Sprintf("%s (ERC1155 #%s)", toAddr, tokenId.String())
		// Append transaction info to the corresponding address+tokenID group
		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		// Track total amount per key
		totals[key] += float64(amount.Int64())
	}

	for addr, txs := range beneficiariesMap {
		results = append(results, Beneficiary{
			BeneficiaryAddress: addr,
			Amount:             totals[addr],
			Transactions:       txs, // All transactions to that address+tokenID
		})
	}

	return results
}
