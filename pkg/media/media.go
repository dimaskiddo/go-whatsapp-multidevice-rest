// Package media provides a unified interface for storing WhatsApp media assets.
// Two backends are supported:
//   - local  – saves files to disk; a built-in HTTP handler serves them back.
//   - s3     – uploads to any S3-compatible store (AWS, MinIO, Cloudflare R2,
//     DigitalOcean Spaces, Backblaze B2, …) using a configurable endpoint.
package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

// Storage is the common interface for all media backends.
type Storage interface {
	// Save persists the raw bytes and returns the public URL that can be
	// embedded in the webhook payload.
	Save(ctx context.Context, filename, mimeType string, data []byte) (string, error)
}

// MediaStorage is the active backend, initialised by init().
var MediaStorage Storage

// StorageType describes which backend is in use ("local", "s3", or "none").
var StorageType string

func init() {
	backend, _ := env.GetEnvString("MEDIA_STORAGE")
	backend = strings.ToLower(strings.TrimSpace(backend))
	StorageType = backend

	switch backend {
	case "local":
		s, err := newLocalStorage()
		if err != nil {
			log.Print(nil).Fatal("Error initializing local media storage: " + err.Error())
		}
		MediaStorage = s
		log.Print(nil).Info("Media storage backend: local (path=" + s.dir + ")")

	case "s3":
		s, err := newS3Storage()
		if err != nil {
			log.Print(nil).Fatal("Error initializing S3 media storage: " + err.Error())
		}
		MediaStorage = s
		log.Print(nil).Info("Media storage backend: s3 (bucket=" + s.bucket + ")")

	default:
		// "none" or unset — media bytes are embedded as base64 in the payload.
		StorageType = "none"
		log.Print(nil).Info("Media storage backend: none (base64 inline)")
	}
}

// ─── Local storage ────────────────────────────────────────────────────────────

type localStorage struct {
	dir     string
	baseURL string // e.g. "https://api.example.com/media"
}

func newLocalStorage() (*localStorage, error) {
	dir, err := env.GetEnvString("MEDIA_LOCAL_PATH")
	if err != nil {
		dir = "./media"
	}

	baseURL, err := env.GetEnvString("MEDIA_LOCAL_BASE_URL")
	if err != nil {
		return nil, errors.New("MEDIA_LOCAL_BASE_URL is required when MEDIA_STORAGE=local")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create media directory %q: %w", dir, err)
	}

	return &localStorage{
		dir:     dir,
		baseURL: strings.TrimRight(baseURL, "/"),
	}, nil
}

func (l *localStorage) Save(_ context.Context, filename, _ string, data []byte) (string, error) {
	dest := filepath.Join(l.dir, filename)
	if err := os.WriteFile(dest, data, 0644); err != nil {
		return "", fmt.Errorf("local save failed for %q: %w", dest, err)
	}
	return l.baseURL + "/" + filename, nil
}

// LocalDir returns the directory used by the local backend so the HTTP handler
// can serve files with http.FileServer.
func LocalDir() string {
	if s, ok := MediaStorage.(*localStorage); ok {
		return s.dir
	}
	return ""
}

// ─── S3-compatible storage ────────────────────────────────────────────────────

type s3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string // optional CDN / public base URL
	keyPrefix string // optional prefix added to every object key
}

func newS3Storage() (*s3Storage, error) {
	bucket, err := env.GetEnvString("MEDIA_S3_BUCKET")
	if err != nil {
		return nil, errors.New("MEDIA_S3_BUCKET is required when MEDIA_STORAGE=s3")
	}

	region, _ := env.GetEnvString("MEDIA_S3_REGION")
	if region == "" {
		region = "us-east-1" // sensible default; MinIO ignores it
	}

	accessKey, err := env.GetEnvString("MEDIA_S3_ACCESS_KEY")
	if err != nil {
		return nil, errors.New("MEDIA_S3_ACCESS_KEY is required when MEDIA_STORAGE=s3")
	}

	secretKey, err := env.GetEnvString("MEDIA_S3_SECRET_KEY")
	if err != nil {
		return nil, errors.New("MEDIA_S3_SECRET_KEY is required when MEDIA_STORAGE=s3")
	}

	// Optional: custom endpoint enables MinIO, R2, Spaces, B2, etc.
	endpointURL, _ := env.GetEnvString("MEDIA_S3_ENDPOINT")

	// Path-style access is required for MinIO and some other providers.
	// Set MEDIA_S3_PATH_STYLE=true to enable it.
	pathStyleStr, _ := env.GetEnvString("MEDIA_S3_PATH_STYLE")
	pathStyle := strings.ToLower(pathStyleStr) == "true"

	// Optional prefix for every stored object, e.g. "whatsapp-media/"
	keyPrefix, _ := env.GetEnvString("MEDIA_S3_KEY_PREFIX")

	// Optional: override the public base URL (CDN, MinIO public bucket URL…).
	// If omitted, a standard AWS-style URL is constructed.
	publicURL, _ := env.GetEnvString("MEDIA_S3_PUBLIC_URL")

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	}

	// Custom timeout (30 s) via http.Client
	opts = append(opts, config.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}))

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("S3 config error: %w", err)
	}

	s3Opts := []func(*s3.Options){}
	if endpointURL != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
		})
	}
	if pathStyle {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(cfg, s3Opts...)

	return &s3Storage{
		client:    client,
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
		keyPrefix: keyPrefix,
	}, nil
}

func (s *s3Storage) Save(ctx context.Context, filename, mimeType string, data []byte) (string, error) {
	key := s.keyPrefix + filename

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", fmt.Errorf("S3 upload failed for %q: %w", key, err)
	}

	// Build the public URL
	var url string
	if s.publicURL != "" {
		url = s.publicURL + "/" + key
	} else {
		// Standard AWS virtual-hosted-style URL
		url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
	}
	return url, nil
}
