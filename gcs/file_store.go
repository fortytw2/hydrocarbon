package gcs

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type FileStore struct {
	client          *storage.Client
	imageBucketName string
}

func NewFileStore(serviceAccount, imageBucketName string) (*FileStore, error) {
	ctx := context.Background()

	creds, err := google.CredentialsFromJSON(context.TODO(), []byte(serviceAccount), storage.ScopeFullControl)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &FileStore{client: client, imageBucketName: imageBucketName}, nil
}

func (fs *FileStore) Put(fileName string, contents []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(contents)
	if err != nil {
		return "", err
	}
	hash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	contentType := http.DetectContentType(contents)

	// only hash the file for now
	fName := hash
	switch contentType {
	case "image/png":
		fName += ".png"
	case "image/jpeg":
		fName += ".jpeg"
	case "image/gif":
		fName += ".gif"
	case "image/webp":
		fName += ".webp"
	default:
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	obj := fs.client.Bucket(fs.imageBucketName).Object(fName)
	wc := obj.NewWriter(ctx)
	written, err := wc.Write(contents)
	if err != nil {
		return "", err
	}

	if written != len(contents) {
		return "", fmt.Errorf("gcs: wrote %d, should have written %d", written, len(contents))
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.storage.googleapis.com/%s", fs.imageBucketName, fName), nil
}

func (fs *FileStore) Stop() error {
	return fs.client.Close()
}
