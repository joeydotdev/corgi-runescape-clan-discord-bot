package storage

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var MissingCredentialsError error = errors.New("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set")

// S3Client is the S3 client to use
var S3Client *s3.Client

// BucketName is the name of the S3 bucket to use
var BucketName string = "corgi-discord-bot"

const (
	// RegionName is the name of the AWS region to use
	RegionName string = "us-west-1"
)

// InitializeS3 initializes the S3 client
func InitializeS3() error {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		return MissingCredentialsError
	}

	customProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(customProvider),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		return err
	}

	S3Client = s3.NewFromConfig(cfg)
	return nil
}

// ListS3Files lists all files in the S3 bucket
func ListS3Files() error {
	resp, err := S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &BucketName,
	}, func(options *s3.Options) {
		options.Region = RegionName
	})

	if err != nil {
		return err
	}

	for _, obj := range resp.Contents {
		println(*obj.Key)
	}

	return nil
}
