package handler

import "fmt"

func validateRequest(message any) error {
  switch message.(type) {
  case *CreateHashRequest:
    req := message.(*CreateHashRequest)
    if req.Payload == nil {
      return fmt.Errorf("empty payload")
    }
  case *CompareHashRequest:
    req := message.(*CompareHashRequest)
    if req.Payload == nil {
      return fmt.Errorf("empty payload")
    }
    if req.ClaimHash == "" {
      return fmt.Errorf("claim hash is mandatory")
    }
  }
  return nil
}
