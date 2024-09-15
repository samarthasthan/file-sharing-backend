package s3

import (
	"fmt"
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

// DeleteFile deletes a file from the S3 bucket
func (s *S3) DeleteFile(filePath string) error {
	key := filepath.Base(filePath)

	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", key, err)
	}

	// Wait until the file is deleted to ensure consistency
	err = s.S3.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to wait for file deletion %s: %w", key, err)
	}

	return nil
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
