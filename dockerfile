FROM golang:1.22.3

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary
RUN go build -o eth-analyzer ./cmd/main.go

EXPOSE 8080

CMD ["./eth-analyzer"]
