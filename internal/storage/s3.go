package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"

	appconfig "github.com/dlsu-lscs/lscs-core-api/internal/config"
)

// S3Config holds S3/Garage configuration
type S3Config struct {
	Endpoint     string
	Bucket       string
	AccessKey    string
	SecretKey    string
	Region       string
	UsePathStyle bool // true for Garage, false for AWS S3
}

// S3Service handles S3 operations
type S3Service struct {
	client *s3.Client
	config S3Config
}

// NewS3Service creates a new S3 service
func NewS3Service(cfg *appconfig.Config) (*S3Service, error) {
	if cfg.S3Endpoint == "" || cfg.S3AccessKeyID == "" || cfg.S3SecretAccessKey == "" {
		log.Info().Msg("S3 configuration not complete, storage service disabled")
		return &S3Service{config: S3Config{Bucket: cfg.S3Bucket}}, nil
	}

	s3Config := S3Config{
		Endpoint:     cfg.S3Endpoint,
		Bucket:       cfg.S3Bucket,
		AccessKey:    cfg.S3AccessKeyID,
		SecretKey:    cfg.S3SecretAccessKey,
		Region:       cfg.S3Region,
		UsePathStyle: true, // Garage requires path-style addressing
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               cfg.S3Endpoint,
			SigningRegion:     cfg.S3Region,
			HostnameImmutable: true,
		}, nil
	})

	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKeyID,
			cfg.S3SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = s3Config.UsePathStyle
	})

	return &S3Service{
		client: client,
		config: s3Config,
	}, nil
}

// IsEnabled returns true if S3 service is configured
func (s *S3Service) IsEnabled() bool {
	return s.client != nil && s.config.Endpoint != ""
}

// GenerateUploadURL generates a pre-signed URL for uploading a profile image
func (s *S3Service) GenerateUploadURL(ctx context.Context, memberID int32, contentType string) (string, string, error) {
	if !s.IsEnabled() {
		return "", "", fmt.Errorf("S3 service not enabled")
	}

	// generate object key: profile-images/{member_id}/{timestamp}.{ext}
	timestamp := time.Now().Format("20060102-150405")
	ext := getExtension(contentType)
	if ext == "" {
		ext = "jpg"
	}
	objectKey := fmt.Sprintf("profile-images/%d/%s.%s", memberID, timestamp, ext)

	// validate content type
	if !isValidImageType(contentType) {
		return "", "", fmt.Errorf("invalid content type: %s (allowed: image/jpeg, image/png, image/webp)", contentType)
	}

	// generate pre-signed URL
	presignClient := s3.NewPresignClient(s.client)
	presignParams := &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}

	presignOutput, err := presignClient.PresignPutObject(ctx, presignParams, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		log.Error().Err(err).Str("object_key", objectKey).Msg("failed to generate presigned URL")
		return "", "", fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return presignOutput.URL, objectKey, nil
}

// GenerateDownloadURL generates a pre-signed URL for downloading a profile image
func (s *S3Service) GenerateDownloadURL(ctx context.Context, objectKey string) (string, error) {
	if !s.IsEnabled() {
		return "", fmt.Errorf("S3 service not enabled")
	}

	presignClient := s3.NewPresignClient(s.client)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectKey),
	}

	presignOutput, err := presignClient.PresignGetObject(ctx, presignParams, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		log.Error().Err(err).Str("object_key", objectKey).Msg("failed to generate download URL")
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return presignOutput.URL, nil
}

// DeleteObject deletes an object from S3
func (s *S3Service) DeleteObject(ctx context.Context, objectKey string) error {
	if !s.IsEnabled() {
		return fmt.Errorf("S3 service not enabled")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Error().Err(err).Str("object_key", objectKey).Msg("failed to delete object")
		return fmt.Errorf("failed to delete object: %w", err)
	}

	log.Debug().Str("object_key", objectKey).Msg("deleted object from S3")
	return nil
}

// ObjectExists checks if an object exists in S3
func (s *S3Service) ObjectExists(ctx context.Context, objectKey string) (bool, error) {
	if !s.IsEnabled() {
		return false, fmt.Errorf("S3 service not enabled")
	}

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

// GetPublicURL returns the public URL for an object
func (s *S3Service) GetPublicURL(objectKey string) string {
	if s.config.Endpoint == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", s.config.Endpoint, s.config.Bucket, objectKey)
}

// helper functions

func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	default:
		return ""
	}
}

func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// ProfileImageResult contains the result of a profile image upload
type ProfileImageResult struct {
	UploadURL   string `json:"upload_url"`
	ObjectKey   string `json:"object_key"`
	DownloadURL string `json:"download_url,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}
