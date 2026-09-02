package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GarageClient struct {
	S3		*s3.Client
	Bucket 	string
} 

func NewGarageClient(ctx context.Context) (*GarageClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("garage"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading aws config %v", err)
		// even it is garage let's call it aws cuz its tuff :)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://s3-api.garage.local")
		// i can change the endpoint when this really goes to aws
		o.UsePathStyle = true
	})


	return &GarageClient{
		S3: client,
		Bucket: "showcase",
	}, nil
}
