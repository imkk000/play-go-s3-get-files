package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	region   = "sgp1"
	bucket   = ""
	endpoint = "https://sgp1.digitaloceanspaces.com"
	key      = ""
	secret   = ""
)

func main() {
	ctx := context.Background()
	cfg, _ := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint}, nil
				})),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")),
	)
	c := s3.NewFromConfig(cfg)
	result, _ := c.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(""),
	})
	for _, o := range result.Contents {
		result, _ := c.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    o.Key,
		})
		fmt.Printf("download %s: ", *o.Key)
		dirname := filepath.Join("data", filepath.Dir(*o.Key))
		if err := os.MkdirAll(dirname, 0755); err != nil {
			if !os.IsExist(err) {
				fmt.Printf("err(%s)\n", err.Error())
				continue
			}
		}
		filename := filepath.Join("data", *o.Key)
		fs, err := os.Create(filename)
		if err != nil {
			fmt.Printf("err(%s)\n", err.Error())
			continue
		}
		fmt.Print(io.Copy(fs, result.Body))
		fmt.Println(fs.Close())
	}
}
