package uploader

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/storage"
)

type ObjectUploader struct {
	bucket *storage.BucketHandle
	path   string
}

func NewObjectUploader(ctx context.Context, bucketName string, path string) (*ObjectUploader, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to GCS: %s", err)
	}

	return &ObjectUploader{
		bucket: client.Bucket(bucketName),
		path:   path,
	}, nil
}

func (o *ObjectUploader) UploadJson(b []byte, objectName string) error {

	objectPath := filepath.Join(o.path, objectName)
	ctx := context.Background()
	w := o.bucket.Object(objectPath).NewWriter(ctx)
	w.ContentType = "application/json"
	_, err := w.Write(b)
	if err != nil {
		return fmt.Errorf("Cannot upload %s: %s", objectName, err)
	}
	defer w.Close()

	return nil
}
