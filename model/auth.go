package model

import (
	"context"
	"errors"
)

var userIdKey = struct{}{}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, &userIdKey, userID)
}

func GetUserID(ctx context.Context) (int, error) {
	v := ctx.Value(&userIdKey)
	userID, ok := v.(int)
	if !ok {
		return 0, errors.New("unset userID")
	}
	return userID, nil
}
