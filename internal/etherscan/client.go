package etherscan

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// FetchTxs handles normal, internal, ERC-20, and ERC-721 transactions
func FetchTxs(address string, action string) ([]byte, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=%s&address=%s&startblock=0&endblock=99999999&sort=desc&apikey=%s",
		action, address, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// FetchERC1155Logs fetches TransferSingle logs from ERC-1155 contracts
func FetchERC1155Logs(contract string, fromBlock string, toBlock string, topic0 string) ([]byte, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=logs&action=getLogs&fromBlock=%s&toBlock=%s&address=%s&topic0=%s&apikey=%s",
		fromBlock, toBlock, contract, topic0, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
