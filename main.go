package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Upstash Redis client
	redisURL := getEnv("UPSTASH_REDIS_REST_URL", "https://moral-pelican-13896.upstash.io")
	redisToken := getEnv("UPSTASH_REDIS_REST_TOKEN", "")

	// Parse Upstash URL to get host and port
	// Format: https://host.upstash.io -> host.upstash.io:6379
	var redisAddr string
	if len(redisURL) > 8 {
		redisAddr = redisURL[8:] + ":6379" // Remove https:// and add port
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      redisAddr,
		Password:  redisToken,
		TLSConfig: &tls.Config{},
	})
	ctx := context.Background()

	// HTTP client
	client := &http.Client{Timeout: 5 * time.Second}

	for {
		// Make HTTP request
		resp, err := client.Get("http://172.20.96.1:8085/data.json")
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Read response
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Failed to read response: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Parse JSON
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Add timestamp
		data["timestamp"] = time.Now().Unix()

		// Convert back to JSON
		jsonData, _ := json.Marshal(data)

		// Add to Redis Stream
		err = rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "hardware:metrics",
			MaxLen: 1000,  // Keep last 1000 readings
			Approx: true,  // Use approximate trimming for better performance
			Values: map[string]interface{}{
				"data": string(jsonData),
			},
		}).Err()
		if err != nil {
			log.Printf("Redis stream add failed: %v", err)
		} else {
			fmt.Printf("Added to stream: %s\n", jsonData)
		}

		time.Sleep(1 * time.Second)
	}
}