package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Order struct {
	ID       string `json:"id"`
	Customer string `json:"customer"`
	Amount   int    `json:"amount"`
}

func main() {
	const (
		numRequests = 50
		concurrency = 10
		baseURL     = "http://localhost:8080"
	)
	
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	
	start := time.Now()
	
	// Create orders
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			
			order := Order{
				ID:       fmt.Sprintf("order-%d", id),
				Customer: fmt.Sprintf("customer-%d", id),
				Amount:   100 + id,
			}
			
			data, _ := json.Marshal(order)
			resp, err := http.Post(baseURL+"/orders", "application/json", bytes.NewBuffer(data))
			if err != nil {
				fmt.Printf("Error creating order %d: %v\n", id, err)
				return
			}
			resp.Body.Close()
		}(i)
	}
	
	wg.Wait()
	
	// Process orders
	resp, err := http.Get(baseURL + "/process")
	if err != nil {
		fmt.Printf("Error processing orders: %v\n", err)
		return
	}
	resp.Body.Close()
	
	fmt.Printf("Load test completed in %v\n", time.Since(start))
}
