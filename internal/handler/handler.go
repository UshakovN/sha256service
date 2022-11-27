package handler

import (
  "context"
  "fmt"
  "sha256service/internal/httpclient"
  "time"
  "crypto/sha256"
)

type Handler struct {
  ctx context.Context
  //hashClient *sha256.SHA256
  httpClient *httpclient.HttpClient
}

func NewRequestHandler(config *Config) *Handler {
  return &Handler{
    //hashClient: config.HashClient,
    httpClient: config.HttpClient,
  }
}

func (h *Handler) GetHash(payload []byte, secret string) (*CreateHashResponse, error) {
  startTime := time.Now()
  //hashSum := h.hashClient.Sum(payload)
  hashSum := sha256.Sum256(payload)

  hashingTime := time.Since(startTime)
  return &CreateHashResponse{
    HashSum:     fmt.Sprintf("%x", hashSum),
    HashedAt:    time.Now().UTC().Format(time.UnixDate),
    HashingTime: hashingTime.String(),
  }, nil
}

func (h *Handler) CompareHash(payload []byte, claimHash string, secret string) (*CompareHashResponse, error) {
  var equal bool
  //hashSum := h.hashClient.Sum(payload)
  hashSum := sha256.Sum256(payload)

  if claimHash == fmt.Sprintf("%x", hashSum) {
    equal = true
  }
  return &CompareHashResponse{
    Equal:      equal,
    ComparedAt: time.Now().UTC().Format(time.UnixDate),
  }, nil
}
