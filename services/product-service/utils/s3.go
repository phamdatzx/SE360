package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"product-service/config"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// UploadImageToS3 uploads an image file to S3 and returns the URL
func UploadImageToS3(file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	if config.S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// Read file content
	buffer := make([]byte, fileHeader.Size)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s/%s-%d%s", folder, uuid.New().String(), time.Now().Unix(), ext)

	// Upload to S3
	_, err = config.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(config.S3BucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:         aws.String("public-read"), // Make the file publicly accessible
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		config.S3BucketName,
		getRegionFromClient(),
		filename,
	)

	return url, nil
}

// DeleteImageFromS3 deletes an image from S3
func DeleteImageFromS3(imageURL string) error {
	if config.S3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	// Extract the key from the URL
	// Assuming URL format: https://bucket.s3.region.amazonaws.com/key
	key := extractKeyFromURL(imageURL)
	if key == "" {
		return fmt.Errorf("invalid S3 URL")
	}

	_, err := config.S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// Helper function to get region from S3 client
func getRegionFromClient() string {
	if config.S3Client != nil && config.S3Client.Config.Region != nil {
		return *config.S3Client.Config.Region
	}
	return "us-east-1" // default region
}

// Helper function to extract key from S3 URL
func extractKeyFromURL(url string) string {
	// Simple extraction - you might want to make this more robust
	// This assumes URL format: https://bucket.s3.region.amazonaws.com/key
	// or https://bucket.s3.amazonaws.com/key
	
	// Find the position after ".amazonaws.com/"
	const suffix = ".amazonaws.com/"
	idx := len(url)
	for i := 0; i < len(url)-len(suffix); i++ {
		if url[i:i+len(suffix)] == suffix {
			return url[i+len(suffix):]
		}
	}
	return ""
}
