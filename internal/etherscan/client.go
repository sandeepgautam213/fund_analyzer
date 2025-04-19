package etherscan

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchTxs(address string, action string) ([]byte, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=%s&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", action, address, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
