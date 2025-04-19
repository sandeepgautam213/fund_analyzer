package analyzer

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
