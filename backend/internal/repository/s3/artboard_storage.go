package s3

import (
	"bytes"
	"fmt"
	"goP2Pbackend/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type artboardStorage struct {
	s3Client *s3.S3
	bucket   string
}

func NewArtboardStorage(s3Client *s3.S3, bucket string) domain.ArtboardStorage {
	return &artboardStorage{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (s *artboardStorage) Save(artboardID string, data []byte) error {
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(artboardID),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *artboardStorage) Load(artboardID string) ([]byte, error) {
	result, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(artboardID),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	return buf.Bytes(), nil
}
