package client

import (
	"app/internal/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type MinioClient struct {
	S3Client *s3.S3
	S3Bucket string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) (*MinioClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(cfg.MinioEndpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(cfg.MinioAccessKeyID, cfg.MinioSecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Error("Failed to create session", slog.String("error", err.Error()))
	}

	// Создание клиента S3
	svc := s3.New(sess)

	// Название существующего бакета
	bucket := cfg.Minio.MinioBucketName

	return &MinioClient{S3Client: svc, S3Bucket: bucket}, nil
}

func (m *MinioClient) UploadFile(uuid uuid.UUID, path string) error {
	const op = "minio.client.UploadFile"

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()
	defer deleteFileAfterUpload(path)

	dir := "upload/"

	_, err = m.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.S3Bucket),
		Key:    aws.String(dir + uuid.String() + ".txt"),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *MinioClient) ReadFile(uid string) (string, error) {
	const op = "minio.client.ReadFile"

	fileKey := fmt.Sprintf("upload/%s.txt", uid)

	result, err := m.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(m.S3Bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer result.Body.Close()

	// Read the object content
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("could not read object body: %w", err)
	}

	return string(body), nil
}

func deleteFileAfterUpload(path string) {
	abs, _ := filepath.Abs(path)
	os.Remove(abs)
}
