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

// UploadFiles uploads multiple files to Firebase Storage
func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	filePath := filepath.Join("../", "serviceAccountKey.json")

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	bucketHandle, err := client.Bucket("shashlikuz-7b2ca.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	for _, v := range file.File["file"] {
		id := uuid.New().String()
		imageFile, err := v.Open()
		if err != nil {
			return nil, err
		}
		defer imageFile.Close()

		fileName := v.Filename

		objectHandle := bucketHandle.Object(fileName)
		writer := objectHandle.NewWriter(context.Background())
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		if _, err := io.Copy(writer, imageFile); err != nil {
			return nil, err
		}
		writer.Close()

		// Encode the filename to handle spaces and special characters
		encodedFileName := url.PathEscape(fileName)
		// Corrected URL
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/shashlikuz-7b2ca.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

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

	// Generate a UUID for the token
	id := uuid.New().String()

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		return nil, err
	}

	// Get the bucket handle
	bucketHandle, err := client.Bucket("shashlikuz-7b2ca.appspot.com")
	if err != nil {
		return nil, err
	}

	// Extract only the base name of the file to avoid including the path
	fileName := filepath.Base(file.Name())

	// Upload the file to Firebase Storage directly
	objectHandle := bucketHandle.Object(fileName)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	defer writer.Close()

	// Copy the file data directly to Firebase
	if _, err := io.Copy(writer, file); err != nil {
		return nil, err
	}

	// Encode the file name to handle special characters
	encodedFileName := url.PathEscape(fileName)

	// Generate the public URL for the file
	fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/shashlikuz-7b2ca.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

	// Add the URL to the response
	resp.Url = append(resp.Url, &models.Url{
		Id:  id,
		Url: fileURL,
	})

	// Return the response with the uploaded file's URL
	return &resp, nil
}

// DeleteFile deletes a file from Firebase Storage
func DeleteFile(id string) error {
	// Initialize a context and Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Bucket name and object path to delete
	bucketName := "shashlikuz-7b2ca.appspot.com"
	objectPath := id

	// Delete the object
	err = client.Objects.Delete(bucketName, objectPath).Do()
	if err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", objectPath, bucketName)
	return nil
}
