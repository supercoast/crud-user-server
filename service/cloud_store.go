package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"cloud.google.com/go/storage"
)

type CloudStore struct {
	googleProjectId  string
	googleBucketName string
}

func NewCloudStore(googleProjectId, googleBucketName string) *CloudStore {
	return &CloudStore{
		googleProjectId:  googleProjectId,
		googleBucketName: googleBucketName,
	}
}

func (c *CloudStore) CreateBucket() (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Not able to generate storage client: %w", err)
	}

	bucket := client.Bucket(c.googleBucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, c.googleProjectId, nil); err != nil {
		return "", fmt.Errorf("Failed to create given bucket: %w", err)
	}

	return c.googleBucketName, nil
}

func (c *CloudStore) Save(imageType string, imageData bytes.Buffer) (string, error) {
	imageId, err := uuid.NewRandom()

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Not able to generate storage client: %w", err)
	}

	bucket := client.Bucket(c.googleBucketName)
	if _, err := bucket.Attrs(ctx); err != nil {
		if err := bucket.Create(ctx, c.googleProjectId, nil); err != nil {
			return "", fmt.Errorf("Failed to create given bucket: %w", err)
		}
	}

	wc := client.Bucket(c.googleBucketName).Object(imageId.String()).NewWriter(ctx)
	if _, err = io.Copy(wc, &imageData); err != nil {
		return "", fmt.Errorf("Could not write file to bucket: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Could not close the client: %w", err)
	}

	return imageId.String(), nil

}
