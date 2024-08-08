package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	UploadObject()
}

func UploadObject() {
	// Create environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Creates the configuration of AWS with the explicit region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	// Open the file
	filePath := "files/hello_world.js"
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer f.Close() // Closing the file after the operation is complete

	// Upload to s3
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("backups-uriel"),
		Key:    aws.String("Test/hello_world.js"),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	fmt.Printf("File uploaded to: %s\n", result.Location)
}
