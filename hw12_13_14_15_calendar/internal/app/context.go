package app

import (
	"context"
)

const UserIDKey = "userID"

type key string

func GetContextValue(ctx context.Context, strKey string) any {
	return ctx.Value(key(strKey))
}

func SetContextValue(ctx context.Context, strKey string, value any) context.Context {
	return context.WithValue(ctx, key(strKey), value)
}
