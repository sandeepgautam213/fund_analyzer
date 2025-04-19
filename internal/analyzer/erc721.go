package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

func ParseERC721Transfers(address string) []Beneficiary {
	nftRaw, err := etherscan.FetchTxs(address, "tokennfttx")
	if err != nil {
		log.Println("Error fetching NFT transfers:", err)
		return nil
	}

	var nftData NFTTxResult
	if err := json.Unmarshal(nftRaw, &nftData); err != nil {
		log.Println("Error parsing NFT transfers:", err)
		return nil
	}

	beneficiariesMap := make(map[string][]TxInfo)
	totals := make(map[string]float64)

	for _, tx := range nftData.Result {
		if tx.To == "" || tx.To == address {
			continue
		}

		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

		txinfo := TxInfo{
			TxAmount:      1, // NFT assumed to be unique per transfer
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}
		if txinfo.TxAmount == 0 {
			continue
		}

		key := fmt.Sprintf("%s (NFT: %s #%s)", tx.To, tx.TokenSymbol, tx.TokenID)
		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		totals[key] += 1
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
