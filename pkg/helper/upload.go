package helper

import (
	"context"
	"fmt"
	"food/api/models"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// UploadFiles uploads multiple files to Firebase Storage
func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	// Load Firebase credentials from environment variable
	filePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if filePath == "" {
		return nil, fmt.Errorf("google credentials path not set")
	}

	// Initialize Firebase App
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	// Initialize Firebase Storage Client
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	// Get a reference to the storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	// Loop through uploaded files and upload to Firebase Storage
	for _, v := range file.File["file"] {
		id := uuid.New().String() // Generate a unique token for download
		imageFile, err := v.Open()
		if err != nil {
			return nil, err
		}
		defer imageFile.Close()

		fileName := v.Filename
		objectHandle := bucketHandle.Object(fileName)
		writer := objectHandle.NewWriter(context.Background())

		// Add token metadata for download
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		// Copy the file data to Firebase Storage
		if _, err := io.Copy(writer, imageFile); err != nil {
			return nil, err
		}
		writer.Close()

		// URL encode the filename to handle spaces and special characters
		encodedFileName := url.PathEscape(fileName)

		// Generate the download URL
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		// Append to the response
		resp.Url = append(resp.Url, &models.Url{
			Id:  id,
			Url: fileURL,
		})
	}

	return &resp, nil
}

func DeleteFile(fileName string) error {
	ctx := context.Background()

	// Load Firebase credentials from environment variable
	filePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if filePath == "" {
		return fmt.Errorf("google credentials path not set")
	}

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(filePath))
	if err != nil {
		return fmt.Errorf("firebase App initialization error: %w", err)
	}

	// Initialize Firebase Storage Client
	client, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("firebase Storage client initialization error: %w", err)
	}

	// Get a reference to the storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		return fmt.Errorf("bucket handle error: %w", err)
	}

	// Get the object to delete
	objectHandle := bucketHandle.Object(fileName)
	if err := objectHandle.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}
