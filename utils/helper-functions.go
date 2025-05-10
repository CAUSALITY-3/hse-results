package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

		// Check accepted encodings
		acceptEncoding := c.Get("Accept-Encoding")

		var tryEncodings []struct {
			ext      string
			encoding string
		}

		if strings.Contains(acceptEncoding, "br") {
			tryEncodings = append(tryEncodings, struct {
				ext, encoding string
			}{".br", "br"})
		}
		if strings.Contains(acceptEncoding, "gzip") {
			tryEncodings = append(tryEncodings, struct {
				ext, encoding string
			}{".gz", "gzip"})
		}

		for _, enc := range tryEncodings {
			compressedPath := filePath + enc.ext
			// Check if the compressed file exists
			fmt.Println("Compressed path", compressedPath)
			if _, err := os.Stat(compressedPath); err == nil {
				c.Set("Content-Encoding", enc.encoding)
				c.Type(filepath.Ext(filePath)) // Set original content type
				return c.SendFile(compressedPath, false)
			}
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
	fmt.Println("Successfully Opened ", fileName)
	if err != nil {
		log.Fatal(err)
		jsonFile.Close()
		return nil, err
	}
	defer jsonFile.Close()

	// Read the JSON file content into a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Create a variable to hold the decoded data
	var typeData T

	// Decode the JSON data into the variable
	json.Unmarshal(byteValue, &typeData)
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
	}
}

func RateLimiterMiddleware(c *fiber.Ctx) error {
	// if rateLimiting {
	ip := c.IP()
	limiter := getVisitor(ip)

	if !limiter.Allow() {
		return c.Status(fiber.StatusTooManyRequests).SendString("Rate limited due to high system load.")
	}
	// }
	return c.Next()
}
