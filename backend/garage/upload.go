package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func (g *GarageClient) UploadProfileImage(ctx context.Context, userID uuid.UUID, r io.Reader, contentType string) (string, error) {
	key := fmt.Sprintf("avatars/%s.jpg", userID.String())

	_, err := g.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(g.Bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("uploading to garage: %w", err)
	}

	return key, nil
}
