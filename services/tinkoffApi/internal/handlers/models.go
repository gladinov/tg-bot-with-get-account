package handlers

type ctxKey string

const (
	EncryptedTokenKey    ctxKey = "X-Encrypted-Token"
	HeaderEncryptedToken        = "X-Encrypted-Token"
)

const (
	ChatIdKey    ctxKey = "X-ChatID"
	HeaderChatID        = "X-ChatID"
)
