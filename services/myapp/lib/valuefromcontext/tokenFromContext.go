package valuefromcontext

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	contextkeys "github.com/gladinov/contracts/context"
)

func GetToken(ctx context.Context) (string, error) {
	tokenBase64, exist := ctx.Value(contextkeys.EncryptedTokenKey).(string)
	if !exist {
		return "", errors.New("context has not token or token not string")
	}
	return tokenBase64, nil
}

func GetChatIDFromCtxStr(ctx context.Context) (string, error) {
	chatID, exist := ctx.Value(contextkeys.ChatIDKey).(string)
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
