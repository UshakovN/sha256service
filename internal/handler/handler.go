package handler

import (
	"context"
	"fmt"
	"net/http"
	"sha256service/internal/httpclient"
	"sha256service/pkg/sha256"
	"time"
)

type Handler struct {
	ctx        context.Context
	hashClient *sha256.SHA256
	httpClient *httpclient.HttpClient
}

func NewRequestHandler(config *Config) *Handler {
	return &Handler{
		hashClient: config.HashClient,
		httpClient: config.HttpClient,
	}
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
