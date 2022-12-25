package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
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

// init initializes the S3 client
func init() {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		panic(MissingCredentialsError)
	}

	customProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(customProvider),
		config.WithRegion(RegionName),
	)
	if err != nil {
		panic(err)
	}

	S3Client = s3.NewFromConfig(cfg)
}

// ListS3Files lists all files in the S3 bucket
func ListS3Files() error {
	resp, err := S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
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

// UploadJSON uploads a JSON blob to S3
func UploadJSON(filename string, data interface{}) error {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(serializedData),
	}, func(options *s3.Options) {
		options.Region = RegionName
	})

	if err != nil {
		return err
	}

	return nil
}

// DownloadJSON downloads a JSON blob from S3
func DownloadJSON(filename string, data interface{}) error {
	resp, err := S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	}, func(options *s3.Options) {
		options.Region = RegionName
	})

	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func ListObjects(prefix string) ([]string, error) {
	resp, err := S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
		Prefix: aws.String(prefix),
	}, func(options *s3.Options) {
		options.Region = RegionName
	})

	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range resp.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}
