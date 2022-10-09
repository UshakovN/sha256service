package handler

import (
	"context"
	"fmt"
	"net/http"
	"sha256service/internal/dynamodb"
	"sha256service/internal/httpclient"
	"sha256service/pkg/sha256"
	"time"
)

type Handler struct {
	ctx            context.Context
	hashClient     *sha256.SHA256
	dynamodbClient *dynamodb.Client
	httpClient     *httpclient.HttpClient
}

func NewRequestHandler(config *Config) *Handler {
	return &Handler{
		hashClient:     config.HashClient,
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

func (h *Handler) GetItemHashFromStore(sum string) (*ItemHash, bool, error) {
	item, found, err := h.dynamodbClient.GetItem(&itemHashKey{
		HashSum: sum,
	}, &ItemHash{})
	if err != nil {
		return nil, false, fmt.Errorf("cannot get item hash from storage %v", err)
	}
	return item.(*ItemHash), found, nil
}

func (h *Handler) GetItemHash(data []byte) (*ItemHash, error) {
	startTime := time.Now()
	hashSum := h.hashClient.Sum(data)
	hashingTime := time.Since(startTime)
	contentType := http.DetectContentType(data)
	return &ItemHash{
		HashFound:   true,
		HashSum:     fmt.Sprintf("%x", hashSum),
		MimeType:    contentType,
		HashedAt:    time.Now().UTC().Format(time.UnixDate),
		HashingTime: hashingTime.String(),
	}, nil
}

func (h *Handler) GetUserSessionFromStore(sessionId string) (*UserSession, bool, error) {
	session, found, err := h.dynamodbClient.GetItem(userSessionKey{
		Id: sessionId,
	}, &UserSession{})
	if err != nil {
		return nil, false, fmt.Errorf("cannot get user session from storage")
	}
	return session.(*UserSession), found, nil
}

func (h *Handler) PutUserSessionToStore(session *UserSession) error {
	if err := h.dynamodbClient.PutItem(session); err != nil {
		return fmt.Errorf("cannot put user session to storage: %v", err)
	}
	return nil
}

func (h *Handler) GetUserFromStore(login string) (*User, bool, error) {
	user, found, err := h.dynamodbClient.GetItem(&userKey{
		Login: login,
	}, &User{})
	if err != nil {
		return nil, false, fmt.Errorf("cannot get user from storage: %v", err)
	}
	return user.(*User), found, nil
}

func (h *Handler) PutUserToStore(user *User) error {
	if err := h.dynamodbClient.PutItem(user); err != nil {
		return fmt.Errorf("cannot put user to storage: %v", err)
	}
	return nil
}
