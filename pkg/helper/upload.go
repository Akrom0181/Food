package helper

import (
	"context"
	"encoding/base64"
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

// UploadFiles uploads multiple files to Firebase Storage and returns their URLs.
func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	// Get Firebase credentials from the environment variable
	credsBase64 := os.Getenv("FIREBASE_CREDENTIALS")
	if credsBase64 == "" {
		log.Println("FIREBASE_CREDENTIALS environment variable is not set")
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS environment variable is not set")
	}

	// Decode the base64-encoded service account key
	credsJson, err := base64.StdEncoding.DecodeString(credsBase64)
	if err != nil {
		log.Println("Failed to decode FIREBASE_CREDENTIALS:", err)
		return nil, err
	}

	// Initialize Firebase app with the credentials
	opt := option.WithCredentialsJSON(credsJson)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	// Upload files
	for _, v := range file.File["file"] {
		id := uuid.New().String()
		imageFile, err := v.Open()
		if err != nil {
			log.Println("Error opening file:", v.Filename, err)
			return nil, err
		}
		defer imageFile.Close()

		log.Println("Uploading file:", v.Filename)

		objectHandle := bucketHandle.Object(v.Filename)
		writer := objectHandle.NewWriter(context.Background())
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		if _, err := io.Copy(writer, imageFile); err != nil {
			log.Printf("Error copying file %s to Firebase Storage: %v", v.Filename, err)
			return nil, err
		}

		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer for file %s: %v", v.Filename, err)
			return nil, err
		}

		log.Println("File uploaded successfully:", v.Filename)

		encodedFileName := url.PathEscape(v.Filename)
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		resp.Url = append(resp.Url, &models.Url{
			Id:  id,
			Url: fileURL,
		})
	}

	return &resp, nil
}

// DeleteFile deletes a file from Firebase Storage using the file ID.
func DeleteFile(id string) error {
	// Get Firebase credentials from the environment variable
	credsBase64 := os.Getenv("FIREBASE_CREDENTIALS")
	if credsBase64 == "" {
		log.Println("FIREBASE_CREDENTIALS environment variable is not set")
		return fmt.Errorf("FIREBASE_CREDENTIALS environment variable is not set")
	}

	// Decode the base64-encoded service account key
	credsJson, err := base64.StdEncoding.DecodeString(credsBase64)
	if err != nil {
		log.Println("Failed to decode FIREBASE_CREDENTIALS:", err)
		return err
	}

	// Initialize Firebase app with the credentials
	opt := option.WithCredentialsJSON(credsJson)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return err
	}

	// Get the Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return err
	}

	// Get a handle to the Firebase Storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return err
	}

	// Delete the file from Firebase Storage
	objectHandle := bucketHandle.Object(id)
	if err := objectHandle.Delete(context.Background()); err != nil {
		log.Println("Failed to delete object:", err)
		return err
	}

	log.Printf("File with ID %s deleted successfully from Firebase Storage.\n", id)
	return nil
}
