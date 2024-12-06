package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kemboy-svg/investment/helpers"
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
        return helpers.JSONResponse(c, http.StatusInternalServerError, false, "API keys not configured", nil)
    }

    // Fetch deposits from Binance API
    data, err := services.FetchDeposits(apiKey, secret)
    if err != nil {
        // Use server error message if available
        return helpers.JSONResponse(c, http.StatusInternalServerError, false, "Failed to fetch deposits: "+err.Error(), nil)
    }

    log.Printf("Raw Response: %s", string(data))

    var deposits []models.Deposit
    if err := json.Unmarshal(data, &deposits); err != nil {
        log.Printf("Error unmarshalling response: %v", err)
        return helpers.JSONResponse(c, http.StatusInternalServerError, false, "Failed to parse deposits: "+err.Error(), nil)
    }

    var syncedDeposits []models.Deposit
    for _, deposit := range deposits {
        // Check if the deposit already exists
        var existing models.Deposit
        if err := store.Db.Where("tx_id = ?", deposit.TxID).First(&existing).Error; err == nil {
            continue // Skip if it already exists
        }

        // Convert deposit amount to USD
        usdAmount, err := services.ConvertToUSD(deposit.Coin, deposit.Amount)
        if err != nil {
            log.Printf("Error converting amount to USD for TxID %s: %v", deposit.TxID, err)
            continue
        }
        deposit.Amount = fmt.Sprintf("%.2f", usdAmount) // Save the USD amount as a string

        // Save the new deposit
        if err := store.Db.Create(&deposit).Error; err != nil {
            return helpers.JSONResponse(c, http.StatusInternalServerError, false, "Failed to save deposit: "+err.Error(), nil)
        }

        // Add to synced deposits
        syncedDeposits = append(syncedDeposits, deposit)
    }

    // Return the response with synced deposits
    return helpers.JSONResponse(c, http.StatusOK, true, "Data fetched successfully", syncedDeposits)
}


func (DepositController) CheckDepositByTxID(c echo.Context) error {
	txID := c.Param("tx_id")
	if txID == "" {
		return helpers.JSONResponse(c, http.StatusBadRequest, false, "tx_id is required", nil)
	}

	var deposit models.Deposit
	if err := store.Db.Where("tx_id = ?", txID).First(&deposit).Error; err != nil {
		return helpers.JSONResponse(c, http.StatusNotFound, false, "Deposit not found", nil)
	}

	response := map[string]interface{}{
		"coin":   deposit.Coin,
		"amount": deposit.Amount,
	}

	return helpers.JSONResponse(c, http.StatusOK, true, "Deposit details", response)
}


func (DepositController) GetAllDeposits(c echo.Context) error {
	// Declare a slice to hold the deposits
	var deposits []models.Deposit

	// Fetch all deposits from the database
	if err := store.Db.Find(&deposits).Error; err != nil {
		return helpers.JSONResponse(c, http.StatusInternalServerError, false, "Failed to fetch deposits", nil)
	}

	// Return the deposits
	return helpers.JSONResponse(c, http.StatusOK, true, "Deposits fetched successfully", deposits)
}




