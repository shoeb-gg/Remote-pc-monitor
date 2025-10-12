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
	// Upstash Redis credentials
	redisAddr := getEnv("UPSTASH_REDIS_ADDR", "moral-pelican-13896.upstash.io:6379")
	redisPassword := getEnv("UPSTASH_REDIS_PASSWORD", "ATZIAAIncDIxMzcxMTEzZjk0NmE0ZTA2YjRhZDUwYzY0MWEwNTQyNHAyMTM4OTY")

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:      redisAddr,
		Password:  redisPassword,
		Username:  "default",
		TLSConfig: &tls.Config{},
	})
	ctx := context.Background()

	// HTTP client
	client := &http.Client{Timeout: 5 * time.Second}
	hardwareMonitorURL := getEnv("HARDWARE_MONITOR_URL", "http://localhost:8085/data.json")

	for {
		// Make HTTP request
		resp, err := client.Get(hardwareMonitorURL)
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
			MaxLen: 1,     // Keep only the latest reading
			Approx: false, // Exact trimming to ensure only 1 entry
			Values: map[string]interface{}{
				"data": string(jsonData),
			},
		}).Err()
		if err != nil {
			log.Printf("Redis stream add failed: %v", err)
		} else {
			fmt.Printf("Logged successfully to stream ✓ - %s\n", time.Now().Format("2006-01-02 15:04:05"))
		}

		time.Sleep(10 * time.Second)
	}
}