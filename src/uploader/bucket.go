package uploader

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/storage"
)

type UploadObject struct {
	ObjectName  string
	ContentType string
	Binary      []byte
}

type GCSUploader struct {
	bucket    *storage.BucketHandle
	directory string
}

func NewGCSUploader(ctx context.Context, bucketName string, directory string) (*GCSUploader, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to GCS: %s", err)
	}

	return &GCSUploader{
		bucket:    client.Bucket(bucketName),
		directory: directory,
	}, nil
}

func (b *GCSUploader) Upload(object *UploadObject) error {
	objectPath := filepath.Join(b.directory, object.ObjectName)
	ctx := context.Background()
	w := b.bucket.Object(objectPath).NewWriter(ctx)
	w.ContentType = object.ContentType
	_, err := w.Write(object.Binary)
	if err != nil {
		return fmt.Errorf("Cannot upload %s to GCS: %s", object.ObjectName, err)
	}
	defer w.Close()

	return nil
}
