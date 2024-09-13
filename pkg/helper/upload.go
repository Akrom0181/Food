package helper

import (
	"context"
	"food/api/models"

	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

// UploadFiles uploads multiple files to Firebase Storage
func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {
	var resp models.MultipleFileUploadResponse

	filePath := filepath.Join("path_to_your_service_account_key", "serviceAccountKey.json")

	opt := option.WithCredentialsFile(filePath)
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

		encodedFileName := url.PathEscape(fileName)
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		resp.Url = append(resp.Url, &models.Url{
			Id:  id,
			Url: fileURL,
		})
	}

	return &resp, nil
}

func DeleteFile(id string) error {
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("path_to_your_service_account_key/serviceAccountKey.json"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	bucketName := "food-8ceb4.appspot.com"
	objectPath := id

	err = client.Objects.Delete(bucketName, objectPath).Do()
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}
