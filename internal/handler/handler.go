package handler

import (
  "context"
  "fmt"
  "sha256service/pkg/sha256"
  "time"
  "sha256service/internal/request"
)

type Handler struct {
  ctx        context.Context
  hashClient *sha256.SHA256
}

func NewRequestHandler(config *Config) *Handler {
  return &Handler{
    ctx:        config.Ctx,
    hashClient: config.HashClient,
  }
}

func (h *Handler) GetHash(payload []byte, secret string) (*request.CreateHashResponse, error) {
  startTime := time.Now()
  hashSum := h.hashClient.Sum(payload, secret)

  hashingTime := time.Since(startTime)
  return &request.CreateHashResponse{
    HashSum:     fmt.Sprintf("%x", hashSum),
    HashedAt:    time.Now().UTC().Format(time.UnixDate),
    HashingTime: hashingTime.String(),
  }, nil
}

func (h *Handler) CompareHash(payload []byte, claimHash string, secret string) (*request.CompareHashResponse, error) {
  var equal bool
  hashSum := h.hashClient.Sum(payload, secret)
  if claimHash == fmt.Sprintf("%x", hashSum) {
    equal = true
  }
  return &request.CompareHashResponse{
    Equal:      equal,
    ComparedAt: time.Now().UTC().Format(time.UnixDate),
  }, nil
}
