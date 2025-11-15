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
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// HardwareNode represents a node in the hardware monitor tree
type HardwareNode struct {
	Text     string          `json:"Text"`
	Value    string          `json:"Value,omitempty"`
	Children []HardwareNode  `json:"Children,omitempty"`
}

// HardwareData represents the root hardware data
type HardwareData struct {
	Children []HardwareNode `json:"Children"`
}

// MetricConfig defines a metric to extract
type MetricConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Path        []string `json:"path"`
	Unit        string   `json:"unit"`
}

// Config holds all metrics configuration
type Config struct {
	Metrics []MetricConfig `json:"metrics"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// loadConfig loads metrics configuration from JSON file
func loadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// parseFloat extracts float from value like "72.3 °C" or "65.7 W"
func parseFloat(value string, unit string) float64 {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.Replace(cleaned, " "+unit, "", 1)

	var result float64
	fmt.Sscanf(cleaned, "%f", &result)
	return result
}

// findNodeByPattern searches for a node matching any of the patterns (separated by |)
func findNodeByPattern(nodes []HardwareNode, pattern string) *HardwareNode {
	patterns := strings.Split(pattern, "|")
	for i := range nodes {
		for _, p := range patterns {
			if strings.Contains(nodes[i].Text, strings.TrimSpace(p)) {
				return &nodes[i]
			}
		}
	}
	return nil
}

// extractMetricByPath extracts a single metric following the path
func extractMetricByPath(computer HardwareNode, path []string, unit string) float64 {
	if len(path) == 0 {
		return 0
	}

	// Start from computer's children
	currentNodes := computer.Children

	// Navigate through the path
	for i, step := range path {
		node := findNodeByPattern(currentNodes, step)
		if node == nil {
			return 0
		}

		// If this is the last step, extract the value
		if i == len(path)-1 {
			return parseFloat(node.Value, unit)
		}

		// Otherwise, go deeper
		currentNodes = node.Children
	}

	return 0
}

// extractMetrics parses the hardware tree using config and extracts all configured metrics
func extractMetrics(data HardwareData, config *Config) map[string]interface{} {
	result := make(map[string]interface{})
	result["timestamp"] = time.Now().Unix()

	if len(data.Children) == 0 {
		return result
	}

	computer := data.Children[0]
	result["pc_name"] = computer.Text

	// Extract each configured metric
	for _, metric := range config.Metrics {
		value := extractMetricByPath(computer, metric.Path, metric.Unit)
		result[metric.Name] = value
	}

	return result
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load metrics configuration
	config, err := loadConfig("metrics-config.json")
	if err != nil {
		log.Fatalf("Failed to load metrics-config.json: %v", err)
	}
	log.Printf("Loaded %d metrics from config", len(config.Metrics))

	// Upstash Redis credentials
	redisAddr := os.Getenv("UPSTASH_REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("UPSTASH_REDIS_ADDR environment variable is required")
	}

	redisPassword := os.Getenv("UPSTASH_REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatal("UPSTASH_REDIS_PASSWORD environment variable is required")
	}

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
	hardwareMonitorURL := os.Getenv("HARDWARE_MONITOR_URL")
	if hardwareMonitorURL == "" {
		hardwareMonitorURL = "http://localhost:8085/data.json"
		log.Printf("Using default HARDWARE_MONITOR_URL: %s", hardwareMonitorURL)
	}

	const maxRetries = 5
	const retryInterval = 12 * time.Second // 5 retries over ~1 minute (12s * 5 = 60s)

	for {
		var rawData HardwareData
		success := false

		// Retry logic: attempt up to maxRetries times
		for attempt := 1; attempt <= maxRetries; attempt++ {
			// Make HTTP request
			resp, err := client.Get(hardwareMonitorURL)
			if err != nil {
				log.Printf("HTTP request failed (attempt %d/%d): %v", attempt, maxRetries, err)
				if attempt < maxRetries {
					time.Sleep(retryInterval)
					continue
				}
				break
			}

			// Read response
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Failed to read response (attempt %d/%d): %v", attempt, maxRetries, err)
				if attempt < maxRetries {
					time.Sleep(retryInterval)
					continue
				}
				break
			}

			// Parse JSON into structured format
			if err := json.Unmarshal(body, &rawData); err != nil {
				log.Printf("Failed to parse JSON (attempt %d/%d): %v", attempt, maxRetries, err)
				if attempt < maxRetries {
					time.Sleep(retryInterval)
					continue
				}
				break
			}

			// Success!
			success = true
			break
		}

		if success {
			// Extract metrics using config
			metrics := extractMetrics(rawData, config)

			// Convert to JSON
			jsonData, err := json.Marshal(metrics)
			if err != nil {
				log.Printf("Failed to marshal metrics JSON: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			// Store in Redis using simple SET (no expiry)
			err = rdb.Set(ctx, "hardware:metrics", string(jsonData), 0).Err()
			if err != nil {
				log.Printf("Redis SET failed: %v", err)
			} else {
				// Display metrics based on config
				cpuTctl := metrics["cpu_temp_tctl"].(float64)
				cpuCCD1 := metrics["cpu_temp_ccd1"].(float64)
				cpuPower := metrics["cpu_power"].(float64)
				gpuTemp := metrics["gpu_temp"].(float64)
				gpuPower := metrics["gpu_power"].(float64)

				fmt.Printf("Metrics stored ✓ [CPU Tctl: %.1f°C | CCD1: %.1f°C | %.1fW | GPU: %.1f°C/%.1fW] - %s\n",
					cpuTctl, cpuCCD1, cpuPower,
					gpuTemp, gpuPower,
					time.Now().Format("15:04:05"))
			}
		} else {
			log.Printf("Failed to fetch hardware data after %d attempts, waiting before next cycle...", maxRetries)
		}

		time.Sleep(10 * time.Second)
	}
}