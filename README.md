# Ethereum Fund Flow Analyzer

This Go-based microservice analyzes Ethereum fund movements for any wallet address using the Etherscan API. It provides structured JSON responses for inflow (/payer) and outflow (/beneficiary) tracing, including support for ETH, ERC-20, ERC-721, and ERC-1155 tokens.

---

## 🔧 Features

- Analyze outbound (beneficiary) and inbound (payer) fund flows
- Tracks:
  - ✅ Normal & internal ETH transactions
  - ✅ ERC-20 token transfers
  - ✅ ERC-721 NFT transfers
  - ✅ ERC-1155 batch token transfers (via logs)
- Includes date-wise breakdown of transactions
- Clean, Dockerized deployment

---

## 📁 Project Structure

```
eth-flow-analyzer/
├── cmd/                # Main entry point
├── config/             # .env loader
├── internal/
│   ├── analyzer/       # Transaction parsing logic
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

Create a file called .env:

```
ETHERSCAN_API_KEY=your_actual_etherscan_key_here
```

### 3. Run Locally (Go CLI)

```bash
go run ./cmd/main.go
```

### 4. Or Run via Docker

Build and start the container:

```bash
docker compose build
docker compose up
```

API will be available at: `http://localhost:8080`

---

## 🔗 API Endpoints

### GET /beneficiary

Tracks where funds have been sent by a wallet address.

Query:
```
/beneficiary?address=0xYourAddressHere
```

Returns:
- List of beneficiary addresses
- Total amount sent to each
- Transaction history with timestamps and tx hashes

### GET /payer

Tracks sources of funds sent to the given address.

Query:
```
/payer?address=0xYourAddressHere
```

Returns:
- List of payer addresses
- Total amount received from each
- Transaction history
- Earliest transaction timestamp ("date")

---

## 📦 Sample Output Format

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

## 📹 Demo Video

🎥 Link: [Google Drive Demo](https://your-demo-link.com)

---

## 🛡 Notes

- Uses public Etherscan API endpoints
- Basic validation for Ethereum address format
- ERC-1155 parsed via log decoding using TransferSingle
- Production-ready Dockerfile + docker-compose

---

## 👨‍💻 Author

Sandeep Gautam  
Blockchain Developer Assignment 

