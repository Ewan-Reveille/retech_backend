package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/google/uuid"
)

const (
    blockchainURL      = "http://localhost:3030"
    productEndpoint    = "/product"
    exchangeEndpoint   = "/exchange"
    httpClientTimeout  = 5 * time.Second
)

var httpClient = &http.Client{Timeout: httpClientTimeout}

type OnChainProduct struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    SellerID    uuid.UUID `json:"seller_id"`
}

type ExchangeEvent struct {
    ProductID uuid.UUID `json:"product_id"`
    From      uuid.UUID `json:"from"`
    To        uuid.UUID `json:"to"`
    Price     float64   `json:"price"`
    Method    string    `json:"method"`
}

func SendProduct(p OnChainProduct) error {
    body, _ := json.Marshal(p)
    resp, err := httpClient.Post(blockchainURL+productEndpoint,
        "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("blockchain error: %s", resp.Status)
    }
    return nil
}

func SendExchange(e ExchangeEvent) error {
    body, _ := json.Marshal(e)
    resp, err := httpClient.Post(blockchainURL+exchangeEndpoint,
        "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("blockchain error: %s", resp.Status)
    }
    return nil
}
