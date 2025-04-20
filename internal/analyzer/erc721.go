package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sandeepgautam213/eth-flow-analyzer/internal/etherscan"
)

// ParseERC721Transfers fetches and parses all ERC-721 (NFT) transfers for a given address.
// It returns a list of Beneficiaries, each containing the NFT receiver, transfer count, and tx history.

func ParseERC721Transfers(address string) []Beneficiary {
	// Fetch raw JSON data of NFT (ERC-721) transfers from Etherscan API
	nftRaw, err := etherscan.FetchTxs(address, "tokennfttx")
	if err != nil {
		log.Println("Error fetching NFT transfers:", err)
		return nil
	}
	// Create a struct to hold parsed NFT transfer data
	var nftData NFTTxResult
	// Unmarshal the raw JSON response into the nftData struct
	if err := json.Unmarshal(nftRaw, &nftData); err != nil {
		log.Println("Error parsing NFT transfers:", err)
		return nil
	}
	// Map to group NFT transfers by recipient and token+ID (since each NFT is unique)
	beneficiariesMap := make(map[string][]TxInfo)
	// Map to keep count of how many NFTs were sent to each recipient
	totals := make(map[string]float64)
	// Loop through each NFT transfer
	for _, tx := range nftData.Result {
		// Skip if there's no recipient or it's a self-transfer
		if tx.To == "" || tx.To == address {
			continue
		}
		// Convert the timestamp string to int64
		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		// Format the timestamp into a readable format
		formattedTime := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
		// Each NFT transfer is considered a single unit (amount = 1)

		txinfo := TxInfo{
			TxAmount:      1, // NFT assumed to be unique per transfer
			DateTime:      formattedTime,
			TransactionID: tx.Hash,
		}
		// Extra safeguard: skip if TxAmount is 0 (though not expected)
		if txinfo.TxAmount == 0 {
			continue
		}
		// Key = recipient + NFT symbol + Token ID to distinguish each NFT

		key := fmt.Sprintf("%s (NFT: %s #%s)", tx.To, tx.TokenSymbol, tx.TokenID)
		// Append this txinfo to the corresponding recipient's record
		beneficiariesMap[key] = append(beneficiariesMap[key], txinfo)
		// Count this NFT transfer
		totals[key] += 1
	}
	// Final result: a slice of Beneficiary structs

	var result []Beneficiary
	// For each recipient key, create a Beneficiary entry
	for addr, txs := range beneficiariesMap {
		result = append(result, Beneficiary{
			BeneficiaryAddress: addr,         // Recipient + NFT info
			Amount:             totals[addr], // Count of NFTs received
			Transactions:       txs,          // List of transfer details
		})
	}

	//List of all beneficiaries with NFT transfer info

	return result
}
