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

type UploadObject struct {
	ObjectName  string
	ContentType string
	Binary      []byte
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

func (o *ObjectUploader) Upload(object *UploadObject) error {

	objectPath := filepath.Join(o.path, object.ObjectName)
	ctx := context.Background()
	w := o.bucket.Object(objectPath).NewWriter(ctx)
	w.ContentType = object.ContentType
	_, err := w.Write(object.Binary)
	if err != nil {
		return fmt.Errorf("Cannot upload %s: %s", object.ObjectName, err)
	}
	defer w.Close()

	return nil
}
