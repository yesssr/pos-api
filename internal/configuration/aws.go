package configuration

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewAwsClient() (*s3.Client, error) {
	accountId := os.Getenv("R2_ACCOUNT_ID");
	accessKey := os.Getenv("ACCESS_KEY");
	secretAccess := os.Getenv("SECRET_ACCESS");

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretAccess,
				"",
			),
		),
		config.WithRegion("auto"),
	);
	if err != nil {
		return nil, err;
	};

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		 o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	});
	return client, nil;
}
