package configuration

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewAwsClient() (*s3.Client, error) {
	endpoint := os.Getenv("AWS_ENDPOINT_URL_S3");
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID");
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY");

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	);
	if err != nil {
		return nil, err;
	};

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint);
		o.UsePathStyle = false;
		o.Region = "auto";
	});
	return client, nil;
}
