package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sandeepgautam213/eth-flow-analyzer/internal/analyzer"
)

func isValidAddress(addr string) bool {
	// Ethereum addresses are 42 characters long and start with 0x followed by 40 hex chars
	matched, _ := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, addr)
	return matched
}

func GetBeneficiaries(c *gin.Context) {
	address := strings.TrimSpace(c.Query("address"))

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "address is required"})
		return
	}

	if !isValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid Ethereum address"})
		return
	}

	result := analyzer.AnalyzeAddress(address)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}
