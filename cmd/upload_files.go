package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func UploadFiles() {
	// Create environment variables
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
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

	errWalk := filepath.Walk("files", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("error: %v", err)
			return err
		}

		//Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file
		f, err := os.Open(path)
		if err != nil {
			log.Printf("failed to open file: %v", err)
			return err
		}
		defer f.Close()

		// Convert path to use the forward slashes
		path = filepath.ToSlash(strings.TrimPrefix(path, "/"))

		//Upload to s3
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET")),
			Key:    aws.String(path),
			Body:   f,
			ACL:    "public-read",
		})
		if err != nil {
			log.Printf("failed to upload file, %v", err)
			return err
		}

		log.Printf("file uploaded to, %s\n", result.Location)
		return nil
	})

	if errWalk != nil {
		log.Printf("failed to walk the folder: %v", errWalk)
	}
}

func UploadSingleFile() {
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
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String("hello_world.js"),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("file uploaded to, %s\n", result.Location)
}
