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

// parseFloat extracts a float from a value like "72.3 °C" or "65.7 W".
// The bool is false when the value is non-numeric (e.g. "N/A", "-", ""),
// so callers can distinguish a real 0 from a failed parse.
func parseFloat(value string, unit string) (float64, bool) {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.Replace(cleaned, " "+unit, "", 1)

	var result float64
	if n, err := fmt.Sscanf(cleaned, "%f", &result); err != nil || n != 1 {
		return 0, false
	}
	return result, true
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

// extractMetricByPath extracts a single metric following the path.
// The bool is false when any path step doesn't match or the value can't be
// parsed, so a missing sensor is never silently reported as 0.
func extractMetricByPath(computer HardwareNode, path []string, unit string) (float64, bool) {
	if len(path) == 0 {
		return 0, false
	}

	// Start from computer's children
	currentNodes := computer.Children

	// Navigate through the path
	for i, step := range path {
		node := findNodeByPattern(currentNodes, step)
		if node == nil {
			return 0, false
		}

		// If this is the last step, extract the value
		if i == len(path)-1 {
			return parseFloat(node.Value, unit)
		}

		// Otherwise, go deeper
		currentNodes = node.Children
	}

	return 0, false
}

// extractMetrics parses the hardware tree using config and extracts all
// configured metrics. Callers must ensure data.Children is non-empty.
// A metric that can't be resolved is stored as nil (JSON null) and logged,
// so the frontend can show "N/A" instead of a misleading 0.
func extractMetrics(data HardwareData, config *Config) map[string]interface{} {
	result := make(map[string]interface{})
	result["timestamp"] = time.Now().Unix()

	computer := data.Children[0]
	result["pc_name"] = computer.Text

	// Extract each configured metric
	for _, metric := range config.Metrics {
		value, ok := extractMetricByPath(computer, metric.Path, metric.Unit)
		if !ok {
			log.Printf("WARNING: metric %q not found (path %v) — emitting null", metric.Name, metric.Path)
			result[metric.Name] = nil
		} else {
			result[metric.Name] = value
		}
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

			// Reject non-200 responses (e.g. a 404/500 error page) instead of
			// letting the error body flow into the JSON parser.
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				log.Printf("Unexpected HTTP status %d (attempt %d/%d)", resp.StatusCode, attempt, maxRetries)
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
			// Guard against a 200-OK-but-empty tree (e.g. monitor returned {} or
			// null): don't clobber good data in Redis with a metrics-less payload.
			if len(rawData.Children) == 0 {
				log.Printf("Hardware data has no children (empty tree) — skipping Redis write")
				time.Sleep(10 * time.Second)
				continue
			}

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
				// Build the summary line from the config so it never drifts from
				// metrics-config.json and can't panic on a missing/null metric.
				var sb strings.Builder
				for _, m := range config.Metrics {
					if sb.Len() > 0 {
						sb.WriteString(" | ")
					}
					if v, ok := metrics[m.Name].(float64); ok {
						fmt.Fprintf(&sb, "%s: %.1f%s", m.Name, v, m.Unit)
					} else {
						fmt.Fprintf(&sb, "%s: N/A", m.Name)
					}
				}
				log.Printf("Metrics stored ✓ [%s]", sb.String())
			}
		} else {
			log.Printf("Failed to fetch hardware data after %d attempts, waiting before next cycle...", maxRetries)
		}

		time.Sleep(10 * time.Second)
	}
}