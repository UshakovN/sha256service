package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
)

type Config struct {
	ctx   context.Context
	aws   aws.Config
	table *string
}

func LoadDefaultConfig(ctx context.Context) *Config {
	awsConfig, err := config.LoadDefaultConfig(ctx,
		func(options *config.LoadOptions) error {
			options.Region = awsRegion
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		ctx:   ctx,
		aws:   awsConfig,
		table: aws.String(baseTable),
	}
}
