package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sandeepgautam213/eth-flow-analyzer/internal/analyzer"
)

func GetPayers(c *gin.Context) {
	address := strings.TrimSpace(c.Query("address"))

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "address is required"})
		return
	}

	if !isValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid Ethereum address"})
		return
	}

	result := analyzer.AnalyzePayers(address)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}
