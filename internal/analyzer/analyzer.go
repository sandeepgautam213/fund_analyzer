package analyzer

import (
	"log"
)

// AnalyzeAddress concurrently fetches ETH, token, NFT, and ERC1155 transactions using Channels
func AnalyzeAddress(address string) []Beneficiary {
	resultChan := make(chan []Beneficiary, 4)

	go func() { resultChan <- ParseETHTransfers(address) }()
	go func() { resultChan <- ParseERC20Transfers(address) }()
	go func() { resultChan <- ParseERC721Transfers(address) }()
	go func() { resultChan <- FetchAndParseERC1155Logs(address) }()

	var result []Beneficiary
	for i := 0; i < 4; i++ {
		partial := <-resultChan
		result = append(result, partial...)
	}

	log.Printf("AnalyzeAddress (Channel): found %d beneficiary entries", len(result))
	return result
}
