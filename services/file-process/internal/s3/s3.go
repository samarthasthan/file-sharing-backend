package s3

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	*s3.S3
	uploader *s3manager.Uploader
	bucket   string
}

// NewS3 creates a new S3 client
func NewS3(accessKey, secretKey, endpoint, region, bucket string) (*S3, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3{
		S3:       svc,
		uploader: uploader,
		bucket:   bucket,
	}, nil
}

// CreateBucket creates an S3 bucket if it doesn't already exist
func (s *S3) CreateBucket() error {
	_, err := s.S3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		if !isErrorBucketAlreadyExists(err) {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

// UploadFile uploads a file to S3
func (s *S3) UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	key := filepath.Base(filePath)
	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file %s: %w", filePath, err)
	}

	publicURL := fmt.Sprintf("http://%s/%s/%s", s.bucket, key)
	return publicURL, nil
}

func isErrorBucketAlreadyExists(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeBucketAlreadyExists:
			return true
		case s3.ErrCodeBucketAlreadyOwnedByYou:
			return true
		default:
			return false
		}
	}
	return false
}
