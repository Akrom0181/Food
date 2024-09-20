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

func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse
	filePath := filepath.Join(".", "serviceAccountKey.json")

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Firebase App initialization error: %v", err)
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Printf("Firebase Storage client initialization error: %v", err)
		return nil, err
	}

	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Printf("Bucket handle error: %v", err)
		return nil, err
	}

	for _, v := range file.File["file"] {
		imageFile, err := v.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", v.Filename, err)
			return nil, err
		}
		defer imageFile.Close() // Ensure file is closed immediately after processing

		id := uuid.New().String()
		objectHandle := bucketHandle.Object(v.Filename)
		writer := objectHandle.NewWriter(context.Background())
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		if _, err := io.Copy(writer, imageFile); err != nil {
			log.Printf("Error copying file to Firebase Storage %s: %v", v.Filename, err)
			return nil, err
		}

		if err := writer.Close(); err != nil {
			log.Printf("Error closing writer for file %s: %v", v.Filename, err)
			return nil, err
		}

		log.Printf("File uploaded successfully: %s", v.Filename)

		encodedFileName := url.PathEscape(v.Filename)
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		resp.Url = append(resp.Url, &models.Url{Id: id, Url: fileURL})
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
		log.Printf("Firebase App initialization error: %v", err)
		return nil, err
	}

	// Initialize Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Printf("Firebase Storage client initialization error: %v", err)
		return nil, err
	}

	// Specify the Firebase Storage bucket
	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Printf("Bucket handle error: %v", err)
		return nil, err
	}

	// Use the base file name
	fileName := filepath.Base(file.Name())

	// Upload the file to Firebase Storage
	objectHandle := bucketHandle.Object(fileName)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	if _, err := io.Copy(writer, file); err != nil {
		log.Printf("Error copying file to Firebase Storage: %v", err)
		return nil, err
	}
	if err := writer.Close(); err != nil {
		log.Printf("Error closing writer for file %s: %v", fileName, err)
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
func DeleteFile(fileName string) error {
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}

	// Bucket name and object path to delete
	bucketName := "food-8ceb4.appspot.com"

	// Delete the object
	err = client.Objects.Delete(bucketName, fileName).Do()
	if err != nil {
		log.Printf("Failed to delete object %s: %v", fileName, err)
		return err
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", fileName, bucketName)
	return nil
}
