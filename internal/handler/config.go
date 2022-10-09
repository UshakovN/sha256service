package handler

import (
	"sha256service/internal/dynamodb"
	"sha256service/internal/httpclient"
)

type Config struct {
	DynamodbClient *dynamodb.Client
	HttpClient     *httpclient.HttpClient
}
