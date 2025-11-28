package models

type ctxKey string

const EncryptedTokenKey ctxKey = "X-Encrypted-Token"
const HeaderEncryptedToken = "X-Encrypted-Token"

const ChatIdKey ctxKey = "X-ChatID"
const HeaderChatID = "X-ChatID"
