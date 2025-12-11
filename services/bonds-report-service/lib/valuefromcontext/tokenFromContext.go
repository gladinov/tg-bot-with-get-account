package valuefromcontext

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

type ctxKey string

const EncryptedTokenKey ctxKey = "X-Encrypted-Token"
const HeaderEncryptedToken = "X-Encrypted-Token"

const ChatIdKey ctxKey = "X-ChatID"
const HeaderChatID = "X-ChatID"

func GetToken(ctx context.Context) (string, error) {
	tokenBase64, exist := ctx.Value(EncryptedTokenKey).(string)
	if !exist {
		return "", errors.New("context has not token or token not string")
	}
	return tokenBase64, nil
}

func GetChatIDFromCtxStr(ctx context.Context) (string, error) {
	chatID, exist := ctx.Value(ChatIdKey).(string)
	if !exist {
		return "", errors.New("context has not chatId or chatId not string")
	}
	return chatID, nil
}

func GetChatIDFromCtxInt(ctx context.Context) (int, error) {
	const op = "processor.chatIDFromContext"
	chatIDStr, err := GetChatIDFromCtxStr(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	return chatID, nil
}
