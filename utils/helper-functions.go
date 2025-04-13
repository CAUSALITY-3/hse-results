package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// func UpdateUsersCache(user userModel.User) bool {
// 	usersCache := SingletonInjector.Get("usersCache").(map[string]*userModel.User)
// 	usersCache[user.Email] = &user
// 	return SingletonInjector.Update(usersCache, "usersCache")
// }

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

// func GetUserCache(email string) *userModel.User {
// 	usersCache := SingletonInjector.Get("usersCache").(map[string]*userModel.User)
// 	return usersCache[email]
// }

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
