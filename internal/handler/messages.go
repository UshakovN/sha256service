package handler

type ItemHash struct {
	HashFound   bool   `json:"hash_found" dynamodbav:"_"`
	HashSum     string `json:"hash_sum" dynamodbav:"hash_sum"`
	MimeType    string `json:"mime_type" dynamodbav:"mime_type"`
	HashedAt    string `json:"hashed_at" dynamodbav:"hashed_at"`
	HashingTime string `json:"hashing_time" dynamodbav:"hashing_time"`
}

type itemHashKey struct {
	HashSum string `dynamodbav:"hash_sum"`
}

type UserSession struct {
	Id            string `json:"id" dynamodbav:"id"`
	UserLogin     string `json:"user_login" dynamodbav:"user_login"`
	Expired       bool   `json:"expired" dynamodbav:"expired"`
	LastLoginDate string `json:"last_login_date" dynamodb:"last_login_date"`
}

type userSessionKey struct {
	Id string `dynamodbav:"id"`
}

type AuthRequest struct {
	Login    string `json:"login" dynamodbav:"login"`
	Password string `json:"password" dynamodbav:"password"`
}

type AuthResponse struct {
	SessionId string `json:"session_id"`
}

type User struct {
	Login    string `json:"login" dynamodbav:"login"`
	Password string `json:"password" dynamodbav:"password"`
}

type userKey struct {
	Login string `dynamodbav:"login"`
}
