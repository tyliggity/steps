package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/stackpulse/public-steps/common/env"
	"io/ioutil"
	"time"
)

var defaultReadTimeout = 30 * time.Second

func DownloadCtx(ctx context.Context, bucket, obj string) ([]byte, error) {
	return download(ctx, env.ReadTimeout(), bucket, obj)
}

func Download(bucket, obj string) ([]byte, error) {
	return download(context.Background(), env.ReadTimeout(), bucket, obj)
}

func download(ctx context.Context, timeout time.Duration, bucket, obj string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't create storage client: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reader, err := client.Bucket(bucket).Object(obj).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't create reader for object %s in buckets %s: %w", obj, bucket, err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("can't read from obj reader: %w", err)
	}

	return data, nil
}
