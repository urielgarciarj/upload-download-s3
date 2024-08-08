package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func DownloadFiles() {
	// Create environment variables
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Creates the configuration of AWS with the explicit region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Printf("error loading AWS configuration: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	bucket := os.Getenv("AWS_BUCKET")
	prefix := "files/" // Adjust the prefix as needed

	// List objects in the bucket
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	// Getting the list of present objects in the bucket
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), listObjectsInput)
	if err != nil {
		log.Fatalf("Error listing objects in bucket: %v", err)
	}

	// Download each file
	for _, object := range listObjectsOutput.Contents {
		key := *object.Key
		filePath := filepath.Join(".", key)

		// Create the directories in the file path if they don't exist
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatalf("Error creating directories: %v", err)
		}

		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()

		// Download the file
		_, err = downloader.Download(context.TODO(), file, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatalf("Error downloading file: %v", err)
		}

		fmt.Printf("File downloaded to: %s\n", filePath)
	}
}
