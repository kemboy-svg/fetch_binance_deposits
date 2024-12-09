package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.binance.com"

func FetchDeposits(apiKey, secret string) ([]byte, error) {
    endpoint := "/sapi/v1/capital/deposit/hisrec"
    timestamp := getCurrentTimestamp()

    queryString := fmt.Sprintf("timestamp=%d", timestamp)

    // Generate signature
    signature := generateHMACSHA256(secret, queryString)

    fullURL := fmt.Sprintf("%s%s?%s&signature=%s", baseURL, endpoint, queryString, signature)

    log.Printf("Request URL: %s", fullURL)

    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return nil, err
    }

    req.Header.Set("X-MBX-APIKEY", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("HTTP request error: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    log.Printf("Response Status: %s", resp.Status)

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Printf("Failed response body: %s", string(body))
        return nil, fmt.Errorf("failed to fetch deposits: HTTP %d - %s", resp.StatusCode, string(body))
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return nil, err
    }

    log.Println("Deposits fetched successfully.")
    return body, nil
}

func getCurrentTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func generateHMACSHA256(secret, data string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}


func ConvertToUSD(coin, amountStr string) (float64, error) {
    const binanceAPIURL = "https://api.binance.com/api/v3/ticker/price"

    // Fetch the coin's price in USD from Binance
    resp, err := http.Get(fmt.Sprintf("%s?symbol=%sUSDT", binanceAPIURL, strings.ToUpper(coin)))
    if err != nil {
        return 0, fmt.Errorf("failed to fetch price for %s: %w", coin, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("received non-200 response: %s", resp.Status)
    }

    var result struct {
        Price string `json:"price"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, fmt.Errorf("failed to decode price response: %w", err)
    }

    price, err := strconv.ParseFloat(result.Price, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid price format: %w", err)
    }
    amount, err := strconv.ParseFloat(amountStr, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid amount format: %w", err)
    }

    return amount * price, nil
}