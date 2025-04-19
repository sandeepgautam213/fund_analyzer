# 🧾 Ethereum Fund Flow Analyzer

This Go-based microservice analyzes Ethereum fund flows for any wallet address using the Etherscan API. It provides structured JSON responses for inflow (`/payer`) and outflow (`/beneficiary`) tracking, including support for ETH, ERC-20, ERC-721, and ERC-1155 tokens.

---

## 🔧 Features

- Analyze outbound (beneficiary) and inbound (payer) fund flows
- Tracks:
  - ✅ Normal & internal ETH transactions
  - ✅ ERC-20 token transfers
  - ✅ ERC-721 NFT transfers
  - ✅ ERC-1155 batch token transfers (via logs)
- Uses concurrency (goroutines + WaitGroups) for fast data aggregation
- Dockerized for easy deployment

---

## 📁 Project Structure

```
eth-flow-analyzer/
├── cmd/                # Main entry point
├── config/             # .env loader
├── internal/
│   ├── analyzer/       # Transaction parsing logic
│   │   ├── eth.go
│   │   ├── erc20.go
│   │   ├── erc721.go
│   │   ├── erc1155.go
│   │   ├── payer.go
│   │   └── analyzer.go
│   ├── handler/        # API endpoints
│   └── etherscan/      # Etherscan client logic
├── go.mod / go.sum
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md
```

---

## 🚀 Quick Start

### 1. Clone the repo

```bash
git clone https://github.com/sandeepgautam213/eth-flow-analyzer.git
cd eth-flow-analyzer
```

### 2. Add Etherscan API Key

Create a file called `.env`:

  cp .env.example .env

```
ETHERSCAN_API_KEY=your_actual_etherscan_key_here
```

### 3. Run Locally (Go CLI)

```bash
go run ./cmd/main.go
```

### 4. Or Run via Docker

```bash
docker compose build
docker compose up
```

API will be available at: `http://localhost:8080`

---

## 🔗 API Endpoints

### GET `/beneficiary?address=0x...`

Tracks outgoing flows (ETH, ERC-20, NFT, ERC-1155) from the address.

### GET `/payer?address=0x...`

Tracks incoming ETH flows to the address.

---

## 🧪 Sample Output Format

```json
{
  "message": "success",
  "data": [
    {
      "beneficiary_address": "0xabc...",
      "amount": 1.5,
      "date": "2025-03-18 09:51:11",
      "transactions": [
        {
          "tx_amount": 0.5,
          "date_time": "2025-03-18 09:45:35",
          "transaction_id": "0x..."
        }
      ]
    }
  ]
}
```

---

## 🧠 ERC-1155 Analyzer Notes

ERC-1155 transfers are decoded using Etherscan's logs API. The service parses `TransferSingle` events:

```
TransferSingle(address operator, address from, address to, uint256 id, uint256 value)
```

Each log is decoded to extract:
- From / To address
- Token ID
- Value transferred
- Timestamp (block time)

---

## 📹 Demo Video

📎 Link: [Google Drive Demo](https://your-demo-link.com)

---

## 🛡 Notes

- Uses public Etherscan API endpoints
- Basic validation for Ethereum address format
- Safe concurrency with goroutines + WaitGroups
- Production-ready Dockerfile + docker-compose

---

## 👨‍💻 Author

[Sandeep Gautam](https://www.linkedin.com/in/gautams1401/)  
Blockchain Developer Assignment – 
