package handler

import (
	"context"
	"sha256service/internal/dynamodb"
	"sha256service/internal/httpclient"
	"sha256service/pkg/sha256"
)

type Config struct {
	Ctx            context.Context
	HashClient     *sha256.SHA256
	DynamodbClient *dynamodb.Client
	HttpClient     *httpclient.HttpClient
}