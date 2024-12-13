package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

	creds := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  "food-8ceb4",
		"private_key_id":              "dcccfd3613f82a2d4a188afaca37eac99ee8ce5d",
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5KxFh2KIJzKE+\ncf9vUPgQblhy5d7Kbb+xX3oIKW3+Ev0D3M8CxzHmZSMfefQVNOyINECc95aifk0b\niGH+q2lO9kO7E2tz6GHfbPb61BdSPRFA3Z9IJ8W1ykwLPb84RduVWeLwSRqjlmcX\nOT8jrTBa0zG/OAcKaz394mmdEgNB9EiS7PuOLQh8hGtTMML//H4RJzfMr6PJ8rE8\n7d+WRZknjYGW5D1dxxsxWTlOxs1rICaxyNfD2n0TOkzhySXtwKyE+teOyPHbu/E3\nsikkfIgVBFy4p8Btwc2UiUvtv+n9XtRqHGBC/jxMEL9iY4ZJfBw6B4hgHj7YAHVs\nnWUJG+tpAgMBAAECggEAApFcftkJujgPvmAcqjcFvE3uGmh3U/hR/GGLaC88OsCZ\nqKLjGLHuzmKgVOXp4sdmm6bvLx4R2olaLrkP9CpSs151PFQTHTrfl0EZawtjkrzi\nZIi465US8qfuW7OMrM5EiB2Efk+NN9qJ7cevSU00CqEQzgrhXMggRnPhgg+cuEk8\naaD+/96zI8vg03Hf7E/pRwbbHo3ohtc1jMigBsCeQjAMx2PwCH2nElNtNLUATEk6\n5AR5+yFK2EOgucM6r65QNx1gKJG7Z9OCcqxeor/HGaST/DdttXnoQpLR4SmML54O\nrR5r/k6NJiGjtMWhcNEtFrv18gUI06o0lodeulOquQKBgQDdgRIWAh9jWxWOzzRW\n4p8hrk6BnREKSdYIoatpjlPV0IKQH4oGjFEiwnZeQUv4sHUSkoKOKu6qgJ/fd3Kl\nvk10nidPk8WbsPuSsCq74PudofRlIozxfpFzkY+IYxHSeVcTMytnItKwdc2vFXZc\nNQg0bwcRAvegAEftw+Xa4DcL7QKBgQDWAVsGWoPyo/yPDngR7K5FimKNgjDczHmA\nwDRKK6FLHg55QLFGVKvY/CYrDxyPivZs1Bp8mSw9jFHx1fPJTntejxq4HcjNHQUd\npQQTs/cTdMvPvVRWxxxVLnkAn9OZFsONbgdt0UmCuYsUIMutr6Ab1ckczzCdzPYB\nd29Pd8tF7QKBgQCeHYqJh05coCJNZP+Znf+2DTUhNLt7OqW8V5uCqASUNlldBAaF\nEhjA1UulkLrodR28+jSTw3XG5DY7UIrYYXXs7xBkr7l5n+aVGYgHwVwbdAZ/QyCV\nKqItexSYaQ/JzLAplnc/Eg6PxCfk+U8aFwkaVL8Yl6On5UtzIEmt6iuhKQKBgHON\nmdPNbi/HIikwm9652MPN3DcilDW05uqBXfqqolYILbKFHvOl5oCsbgOUDkznsPXE\ndWTP5FZ7fQfDCfapvO2rAbdmxbUTNV7zakclRoUn7KEITxDoREEubcHLixq/cunb\n/oDqn/HJM/KzXqczDJXbEtPOgCbEtBTIo77aJVVlAoGBAMHAak1r+qJKAKgYLP4a\nBCJvMF7YT8wvwdTHuE+YvboASWVEtfZGNlkUdZgbuEXizzKHoOwzD/cuqcZ+5VcE\nOjZ01hrQmfe1k9I2W4XcRQvBvtXndmLrVzD81t/iMFzWtDXxw+vs9Sj+Ve/uPPWr\nz7HU90z4z9sOYC4COBnnee8u\n-----END PRIVATE KEY-----\n",
		"client_email":                "firebase-adminsdk-i8uol@food-8ceb4.iam.gserviceaccount.com",
		"client_id":                   "105015910200347949932",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-i8uol%40food-8ceb4.iam.gserviceaccount.com",
		"universe_domain":             "googleapis.com",
	}

	// Marshal the creds map into a JSON byte array
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		log.Println("Error marshaling Firebase credentials:", err)
		return nil, err
	}

	// Use the credentials to initialize the Firebase app
	opt := option.WithCredentialsJSON([]byte(credsJSON))
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
