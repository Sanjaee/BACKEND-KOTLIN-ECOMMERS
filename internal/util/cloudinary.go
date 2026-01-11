package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"
)

type CloudinaryUploader struct {
	CloudName string
	APIKey    string
	APISecret string
}

func NewCloudinaryUploader(cloudName, apiKey, apiSecret string) *CloudinaryUploader {
	return &CloudinaryUploader{
		CloudName: cloudName,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

// UploadImage uploads a single image to Cloudinary and returns the secure URL
// Uses transformations: w_1080,h_1080,c_limit,q_auto,f_auto for optimization
func (c *CloudinaryUploader) UploadImage(fileData []byte, fileName string, folder string) (string, error) {
	// Generate signature
	timestamp := time.Now().Unix()
	transformation := "w_1080,h_1080,c_limit,q_auto,f_auto" // Optimize: resize, compress, auto format
	signature := c.generateSignatureWithTransformation(timestamp, folder, transformation)

	// Create multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields
	writer.WriteField("api_key", c.APIKey)
	writer.WriteField("timestamp", fmt.Sprintf("%d", timestamp))
	writer.WriteField("signature", signature)
	writer.WriteField("transformation", transformation)
	if folder != "" {
		writer.WriteField("folder", folder)
	}
	writer.WriteField("resource_type", "image")

	// Add file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Make request
	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.CloudName)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cloudinary upload failed: %s (status: %d)", string(body), resp.StatusCode)
	}

	// Parse JSON response
	var response struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if response.SecureURL != "" {
		return response.SecureURL, nil
	}
	return response.URL, nil
}

// UploadMultipleImages uploads multiple images to Cloudinary
func (c *CloudinaryUploader) UploadMultipleImages(files []FileData, folder string, maxImages int) ([]string, error) {
	if len(files) > maxImages {
		return nil, fmt.Errorf("maximum %d images allowed, got %d", maxImages, len(files))
	}

	var urls []string
	for _, file := range files {
		url, err := c.UploadImage(file.Data, file.Name, folder)
		if err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", file.Name, err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

type FileData struct {
	Data []byte
	Name string
}

// generateSignature generates Cloudinary signature (backward compatible)
func (c *CloudinaryUploader) generateSignature(timestamp int64, folder string) string {
	return c.generateSignatureWithTransformation(timestamp, folder, "")
}

// generateSignatureWithTransformation generates Cloudinary signature with transformation
func (c *CloudinaryUploader) generateSignatureWithTransformation(timestamp int64, folder string, transformation string) string {
	params := map[string]string{
		"timestamp": fmt.Sprintf("%d", timestamp),
	}
	if folder != "" {
		params["folder"] = folder
	}
	if transformation != "" {
		params["transformation"] = transformation
	}

	// Sort keys and build string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramString := strings.Join(parts, "&") + c.APISecret

	hash := sha1.Sum([]byte(paramString))
	return fmt.Sprintf("%x", hash)
}
