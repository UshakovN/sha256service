package handler

type ItemHash struct {
	HashSum     string `json:"hash_sum" dynamodbav:"hash_sum"`
	MimeType    string `json:"mime_type" dynamodbav:"mime_type"`
	HashedAt    string `json:"hashed_at" dynamodbav:"hashed_at"`
	HashingTime string `json:"hashing_time" dynamodbav:"hashing_type"`
}

type ItemHashResponse ItemHash

type ItemHashKey struct {
	HashSum string `dynamodbav:"hash_sum"`
}
