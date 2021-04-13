package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gocloud.dev/blob"
)

type storage struct {
	Bucket string `env:"BUCKET,required"`
}

func (s storage) initGCloud() error {
	// Note: our implementation injects the base64 representation of the credentials file from the integration.
	// in order to seamlessly support it using go-cloud we will unpack the credentials to the expected format.
	authCode := os.Getenv("AUTH_CODE")
	if authCode == "" {
		return nil
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(authCode)
	if err != nil {
		return fmt.Errorf("decode auth code: %w", err)
	}

	file, err := ioutil.TempFile("", "gcloud-auth")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer closeWithLog(file)

	_, err = file.Write(decodedAuth)
	if err != nil {
		return fmt.Errorf("write decoded auth to temp file: %w", err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file.Name())
	if err != nil {
		return fmt.Errorf("set GOOGLE_APPLICATION_CREDENTIALS to the credentials file: %w", err)
	}

	return nil
}

func (s storage) uploadFile(ctx context.Context, bucket *blob.Bucket, file, key string) error {
	w, err := bucket.NewWriter(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("new bucket writer for key: '%s': %w", key, err)
	}
	defer closeWithLog(w)

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open file: '%s': %w", file, err)
	}
	defer closeWithLog(f)

	_, err = io.Copy(w, f)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func (s storage) upload(ctx context.Context, files map[string]string) error {
	bucket, err := blob.OpenBucket(ctx, s.Bucket)
	if err != nil {
		return fmt.Errorf("open bucket: '%s': %w", s.Bucket, err)
	}
	defer closeWithLog(bucket)

	for file, key := range files {
		err = s.uploadFile(ctx, bucket, file, key)
		if err != nil {
			return err
		}
	}

	return nil
}
