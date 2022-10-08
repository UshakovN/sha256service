package handler

type HashSumResponse struct {
	HashSum     string `json:"hash_sum" dynamodbav:"hash_sum"`
	MimeType    string `json:"mime_type"`
	HashedAt    string `json:"hashed_at"`
	HashingTime string `json:"hashing_time"`
}
