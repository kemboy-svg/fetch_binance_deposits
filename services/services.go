package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://api.binance.com"

func FetchDeposits(apiKey, secret string) ([]byte, error) {
    endpoint := "/sapi/v1/capital/deposit/hisrec"
    timestamp := getCurrentTimestamp()

    // Prepare query string
    queryString := fmt.Sprintf("timestamp=%d", timestamp)

    // Generate signature
    signature := generateHMACSHA256(secret, queryString)

    // Append signature to query string
    fullURL := fmt.Sprintf("%s%s?%s&signature=%s", baseURL, endpoint, queryString, signature)

    // Log the full URL for debugging
    log.Printf("Request URL: %s", fullURL)

    // Create request
    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        return nil, err
    }

    req.Header.Set("X-MBX-APIKEY", apiKey)

    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("HTTP request error: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Log the response status
    log.Printf("Response Status: %s", resp.Status)

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Printf("Failed response body: %s", string(body))
        return nil, fmt.Errorf("failed to fetch deposits: HTTP %d - %s", resp.StatusCode, string(body))
    }

    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return nil, err
    }

    // Log success
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
