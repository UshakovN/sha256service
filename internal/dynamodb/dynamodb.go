package dynamodb

import (
	"context"
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
		return err
	}
	if _, err = client.aws.PutItem(context.TODO(), &db.PutItemInput{
		TableName: client.table,
		Item:      marshalItem,
	}); err != nil {
		return err
	}
	return nil
}

func (client *Client) GetItem(keyValues interface{}, itemStruct interface{}) (interface{}, error) {
	marshalKeys, err := attributevalue.MarshalMap(keyValues)
	if err != nil {
		return nil, err
	}
	marshalItem, err := client.aws.GetItem(context.TODO(), &db.GetItemInput{
		TableName: client.table,
		Key:       marshalKeys,
	})
	if err != nil {
		return nil, err
	}
	item := itemStruct
	if err = attributevalue.UnmarshalMap(marshalItem.Item, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (client *Client) DeleteItem(keyValues interface{}) error {
	marshalKeys, err := attributevalue.MarshalMap(keyValues)
	if err != nil {
		return err
	}
	if _, err = client.aws.DeleteItem(context.TODO(), &db.DeleteItemInput{
		TableName: client.table,
		Key:       marshalKeys,
	}); err != nil {
		return err
	}
	return nil
}

func (client *Client) WriteBatch(items []interface{}) error {
	requests := make([]types.WriteRequest, 0)
	for _, item := range items {
		marshalItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}
		requests = append(requests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: marshalItem,
			},
		})
	}
	if _, err := client.aws.BatchWriteItem(context.TODO(), &db.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			*client.table: requests,
		},
	},
	); err != nil {
		return err
	}
	return nil
}
