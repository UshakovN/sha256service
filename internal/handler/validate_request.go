package handler

import (
  "fmt"
  "sha256service/internal/request"
)

func validateRequest(message any) error {
  switch message.(type) {
  case *request.CreateHashRequest:
    req := message.(*request.CreateHashRequest)
    if req.Payload == nil {
      return fmt.Errorf("empty payload")
    }
  case *request.CompareHashRequest:
    req := message.(*request.CompareHashRequest)
    if req.Payload == nil {
      return fmt.Errorf("empty payload")
    }
    if req.ClaimHash == "" {
      return fmt.Errorf("claim hash is mandatory")
    }
  }
  return nil
}
