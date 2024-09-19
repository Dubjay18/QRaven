package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get Cloudinary credentials from environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Initialize Cloudinary instance
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Open the file to read its content
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer file.Close()

	// Upload file to Cloudinary
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: "avatars", // Folder in your Cloudinary account where the images will be stored
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %v", err)
	}

	// Return the URL of the uploaded image
	return uploadResult.SecureURL, nil
}
