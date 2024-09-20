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
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

// UploadFiles uploads multiple files to Firebase Storage and returns their URLs.
func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	// Path to your Firebase service account key file
	filePath := filepath.Join(".", "serviceAccountKey.json")

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	// Get a handle to the Firebase Storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	// Iterate over the uploaded files
	for _, v := range file.File["file"] {
		id := uuid.New().String()
		imageFile, err := v.Open()
		if err != nil {
			log.Println("Error opening file:", v.Filename, err)
			return nil, err
		}
		defer imageFile.Close() // Ensure file is closed after processing

		fileName := v.Filename
		log.Println("Uploading file:", fileName)

		// Upload the file to Firebase Storage
		objectHandle := bucketHandle.Object(fileName)
		writer := objectHandle.NewWriter(context.Background())
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		if _, err := io.Copy(writer, imageFile); err != nil {
			log.Println("Error copying file to Firebase Storage:", err)
			return nil, err
		}

		// Closing the writer to complete the upload
		if err := writer.Close(); err != nil {
			log.Println("Error closing writer:", err)
			return nil, err
		}

		log.Println("File uploaded successfully:", fileName)

		// Encode the file name to handle spaces and special characters
		encodedFileName := url.PathEscape(fileName)
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		// Append the file URL to the response
		resp.Url = append(resp.Url, &models.Url{
			Id:  id,
			Url: fileURL,
		})
	}

	return &resp, nil
}

// UploadFile uploads a single file to Firebase Storage
func UploadFile(file *os.File) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	id := uuid.New().String()

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	// Specify the Firebase Storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	// Use the original file name
	fileName := file.Name()

	// Upload the file to Firebase Storage
	objectHandle := bucketHandle.Object(fileName)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	if _, err := io.Copy(writer, file); err != nil {
		log.Println("Error copying file to Firebase Storage:", err)
		return nil, err
	}
	if err := writer.Close(); err != nil {
		log.Println("Error closing writer:", err)
		return nil, err
	}

	// Encode the file name to handle spaces and special characters
	encodedFileName := url.PathEscape(fileName)
	fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

	// Add the URL to the response
	resp.Url = append(resp.Url, &models.Url{
		Id:  id,
		Url: fileURL,
	})

	return &resp, nil
}

// DeleteFile deletes a file from Firebase Storage
func DeleteFile(id string) error {
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Println("Failed to create client:", err)
		return err
	}

	// Bucket name and object path to delete
	bucketName := "food-8ceb4.appspot.com"
	objectPath := id

	// Delete the object
	err = client.Objects.Delete(bucketName, objectPath).Do()
	if err != nil {
		log.Println("Failed to delete object:", err)
		return err
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", objectPath, bucketName)
	return nil
}
