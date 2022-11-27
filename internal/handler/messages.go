package handler

import (
  "fmt"
  "encoding/base64"
)

const (
  PayloadTypeString      = "string"
  PayloadTypeEncodedFile = "encoded-file"
)

type CreateHashRequest struct {
  Secret      string `json:"secret"`
  Payload     any    `json:"payload"`
  PayloadType string `json:"payload_type"`
}

type CreateHashResponse struct {
  HashSum     string `json:"hash_sum"`
  HashedAt    string `json:"hashed_at"`
  HashingTime string `json:"hashing_time"`
}

type CompareHashRequest struct {
  Secret      string `json:"secret"`
  ClaimHash   string `json:"claim_hash"`
  Payload     any    `json:"payload"`
  PayloadType string `json:"payload_type"`
}

type CompareHashResponse struct {
  Equal      bool   `json:"equal"`
  ComparedAt string `json:"compared_at"`
}

func preparePayload(payload any, payloadType string) ([]byte, error) {
  var b []byte
  switch payloadType {

  case PayloadTypeEncodedFile:
    payloadStr, ok := payload.(string)
    if !ok {
      return nil, fmt.Errorf("invalid file base64 string")
    }
    var err error
    b, err = base64.StdEncoding.DecodeString(payloadStr)
    if err != nil {
      return nil, fmt.Errorf("cannot decode file base64 string")
    }

  case PayloadTypeString:
    fallthrough
  default:
    payloadStr, ok := payload.(string)
    if !ok {
      return nil, fmt.Errorf("invalid plain text")
    }
    b = []byte(payloadStr)
  }

  return b, nil
}
