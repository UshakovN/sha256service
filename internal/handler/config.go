package handler

import (
  "context"
  "sha256service/pkg/sha256"
)

type Config struct {
  Ctx        context.Context
  HashClient *sha256.SHA256
}
