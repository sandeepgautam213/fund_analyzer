package analyzer

import (
	"log"
	"sync"
)

type resultMsg struct {
	data []Beneficiary
}

// AnalyzeAddress concurrently fetches ETH, token, NFT, and ERC1155 transactions using WaitGroup
func AnalyzeAddress(address string) []Beneficiary {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var result []Beneficiary

	wg.Add(4)

	go func() {
		defer wg.Done()
		eth := ParseETHTransfers(address)
		mu.Lock()
		result = append(result, eth...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		erc20 := ParseERC20Transfers(address)
		mu.Lock()
		result = append(result, erc20...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		erc721 := ParseERC721Transfers(address)
		mu.Lock()
		result = append(result, erc721...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		erc1155 := FetchAndParseERC1155Logs(address)
		mu.Lock()
		result = append(result, erc1155...)
		mu.Unlock()
	}()

	wg.Wait()

	log.Printf("AnalyzeAddress (WaitGroup): found %d beneficiary entries", len(result))
	return result
}
