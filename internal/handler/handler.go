package handler

import (
	"fmt"
	"sha256service/internal/dynamodb"
	"sha256service/internal/httpclient"
)

type Handler struct {
	dynamodbClient *dynamodb.Client
	httpClient     *httpclient.HttpClient
}

func NewRequestHandler(config *Config) *Handler {
	return &Handler{
		dynamodbClient: config.DynamodbClient,
		httpClient:     config.HttpClient,
	}
}

func (h *Handler) PutItemHashInStore(item *ItemHash) error {
	if err := h.dynamodbClient.PutItem(item); err != nil {
		return fmt.Errorf("cannot put item hash to storage: %v", err)
	}
	return nil
}

func (h *Handler) GetItemHashFromStore(sum string) (*ItemHash, error) {
	itemHash, err := h.dynamodbClient.GetItem(&ItemHashKey{
		HashSum: sum,
	}, &ItemHash{})
	if err != nil {
		return nil, fmt.Errorf("cannot get item hash from storage")
	}
	return itemHash.(*ItemHash), nil
}
