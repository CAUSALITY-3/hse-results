package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/time/rate"
)

func ParseBody[T any](c *fiber.Ctx) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

func ServeCompressedFile(root string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filePath := filepath.Join(root, c.Params("*"))

		if _, err := os.Stat(filePath + ".gz"); err == nil {

			fmt.Println("Got Compressed path", filePath+".gz")
			c.Set("Content-Encoding", "gzip")
			c.Type(filepath.Ext(filePath))
			return c.SendFile(filePath+".gz", false)
		}
		// fallback: serve original file
		return c.SendFile(filePath, false)
	}
}

func Filter[T any](slice []T, condition func(T) bool) []T {
	var result []T = []T{}
	for _, item := range slice {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](slice []T, transform func(T) R) []R {
	var result []R = make([]R, len(slice))
	for index, item := range slice {
		result[index] = transform(item)
	}
	return result
}

func Find[T any](slice []T, condition func(T) bool) *T {
	for _, item := range slice {
		if condition(item) {
			return &item
		}
	}
	return nil
}

func Includes[T any](slice []T, condition func(T) bool) bool {
	log.Println("Slice", slice)
	for _, item := range slice {
		if condition(item) {
			return true
		}
	}
	return false
}

func ReadFile[T any](fileName string) (*T, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		return nil, err
	}
	defer func() {
		if cerr := jsonFile.Close(); cerr != nil {
			log.Printf("Error closing file %s: %v", fileName, cerr)
		}
	}()

	// Read the JSON file content into a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Error reading file %s: %v", fileName, err)
		return nil, err
	}

	// Create a variable to hold the decoded data
	var typeData T

	// Decode the JSON data into the variable
	if err := json.Unmarshal(byteValue, &typeData); err != nil {
		log.Printf("Error unmarshalling JSON from file %s: %v", fileName, err)
		return nil, err
	}

	return &typeData, nil
}

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors      = make(map[string]*Visitor)
	mu            sync.Mutex
	rateLimiting  = true
	cpuThreshold  = 75.0
	memThreshold  = 75.0
	checkInterval = 5 * time.Second
)

func MnitorSystem() {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		memStats, _ := mem.VirtualMemory()

		log.Printf("CPU Usage: %.2f%%, Memory Usage: %.2f%%\n", cpuPercent[0], memStats.UsedPercent)
		if cpuPercent[0] > cpuThreshold || memStats.UsedPercent > memThreshold {
			rateLimiting = true
		} else {
			rateLimiting = false
		}

		time.Sleep(checkInterval)
	}
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Visitor IP", ip)
	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 10) // 5 requests per second, burst of 10
		visitors[ip] = &Visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func CeanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
		statsMu.Lock()
		activeUserStats = append(activeUserStats, ActiveUserStat{
			Timestamp: time.Now().Truncate(time.Second),
			Count:     len(visitors),
		})
		statsMu.Unlock()
	}
}

func getClientIP(c *fiber.Ctx) string {
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := c.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	return c.IP() // fallback
}

func RateLimiterMiddleware(c *fiber.Ctx) error {
	// if rateLimiting {
	// return c.Next()

	ip := getClientIP(c)
	limiter := getVisitor(ip)

	if !limiter.Allow() {
		return c.Status(fiber.StatusTooManyRequests).SendString("Rate limited due to high system load.")
	}
	// }
	return c.Next()
}

func ExtractRollNo(input string) string {
	re := regexp.MustCompile(`\((\d+)\)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

type ActiveUserStat struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

var (
	activeUserStats []ActiveUserStat
	statsMu         sync.Mutex
)

func ActiveUsers() int {
	mu.Lock()
	defer mu.Unlock()
	count := 0
	now := time.Now()
	for _, v := range visitors {
		if now.Sub(v.lastSeen) <= 3*time.Minute {
			count++
		}
	}
	return count
}

func TrackActiveUsers() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		count := len(visitors)
		statsMu.Lock()
		activeUserStats = append(activeUserStats, ActiveUserStat{
			Timestamp: time.Now().Truncate(time.Second),
			Count:     count,
		})
		statsMu.Unlock()
	}
}

func PersistActiveUserStats() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		<-ticker.C
		statsMu.Lock()
		if len(activeUserStats) > 0 {
			filename := fmt.Sprintf("./active-users/active_users_%s.json", time.Now().Format("2006-01-02T15"))
			file, err := os.Create(filename)
			if err == nil {
				enc := json.NewEncoder(file)
				enc.SetIndent("", "  ")
				enc.Encode(activeUserStats)
				file.Close()
			} else {
				log.Printf("Error creating stats file: %v", err)
			}
			activeUserStats = nil
		}
		statsMu.Unlock()
	}
}
