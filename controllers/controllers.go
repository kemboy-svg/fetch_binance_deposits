package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/kemboy-svg/investment/models"
	"github.com/kemboy-svg/investment/services"
	"github.com/kemboy-svg/investment/store"
	"github.com/labstack/echo/v4"
)

type DepositController struct{}

func (DepositController) SyncDeposits(c echo.Context) error {
    apiKey := os.Getenv("API_KEY")
	secret := os.Getenv("SECRET_KEY")

	if apiKey == "" || secret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "API keys not configured"})
	}

    // Fetch deposits from Binance API
    data, err := services.FetchDeposits(apiKey, secret)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch deposits"})
    }

    log.Printf("Raw Response: %s", string(data))

    var deposits []models.Deposit
    if err := json.Unmarshal(data, &deposits); err != nil {
        log.Printf("Error unmarshalling response: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse deposits"})
    }

    for _, deposit := range deposits {
        // Check if the deposit already exists
        var existing models.Deposit
        if err := store.Db.Where("tx_id = ?", deposit.TxID).First(&existing).Error; err == nil {
            continue // Skip if it already exists
        }


        // Save the new deposit
        if err := store.Db.Create(&deposit).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save deposit"})
        }
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Deposits synced successfully"})
}
