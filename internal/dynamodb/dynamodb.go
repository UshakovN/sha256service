package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	db "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	baseTable = "sha256-service"
	awsRegion = "us-east-2"
)

type Client struct {
	ctx   context.Context
	aws   *db.Client
	table *string
}

func NewDynamodbClient(config *Config) *Client {
	return &Client{
		ctx:   config.ctx,
		aws:   db.NewFromConfig(config.aws),
		table: config.table,
	}
}

func (client *Client) PutItem(item interface{}) error {
	marshalItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("cannot marshal item to dynamodbav: %v", err)
	}
	if _, err = client.aws.PutItem(client.ctx, &db.PutItemInput{
		TableName: client.table,
		Item:      marshalItem,
	}); err != nil {
		return fmt.Errorf("cannot put item to dynamodb: %v", err)
	}
	return nil
}

func (client *Client) GetItem(keyValues interface{}, itemStruct interface{}) (interface{}, bool, error) {
	marshalKeys, err := attributevalue.MarshalMap(keyValues)
	if err != nil {
		return nil, false, fmt.Errorf("cannot marshal item to dynamodbav: %v", err)
	}
	marshalItem, err := client.aws.GetItem(client.ctx, &db.GetItemInput{
		TableName: client.table,
		Key:       marshalKeys,
	})
	if err != nil {
		return nil, false, fmt.Errorf("cannot get item from dynamodb: %v", err)
	}
	item := itemStruct
	if marshalItem.Item == nil {
		return item, false, nil
	}
	if err = attributevalue.UnmarshalMap(marshalItem.Item, item); err != nil {
		return nil, false, fmt.Errorf("cannot unmarshal dynamodbav item: %v", err)
	}
	return item, true, nil
}

func (client *Client) DeleteItem(keyValues interface{}) error {
	marshalKeys, err := attributevalue.MarshalMap(keyValues)
	if err != nil {
		return fmt.Errorf("cannot marshal item to dynamodbav: %v", err)
	}
	if _, err = client.aws.DeleteItem(client.ctx, &db.DeleteItemInput{
		TableName: client.table,
		Key:       marshalKeys,
	}); err != nil {
		return fmt.Errorf("cannot delete item from dynamodb: %v", err)
	}
	return nil
}

func (client *Client) WriteBatch(items []interface{}) error {
	requests := make([]types.WriteRequest, 0)
	for _, item := range items {
		marshalItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			return fmt.Errorf("cannot marshal item to dynamodbav: %v", err)
		}
		requests = append(requests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: marshalItem,
			},
		})
	}
	if _, err := client.aws.BatchWriteItem(client.ctx, &db.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			*client.table: requests,
		},
	},
	); err != nil {
		return fmt.Errorf("cannot write items batch to dynamodb: %v", err)
	}
	return nil
}
