package model

import (
	"context"
	"errors"
)

type requestUserInfo struct {
	token  string
	userID int
}

var userInfoKey = struct{}{}

func SetUserInfo(ctx context.Context, token string, userID int) context.Context {
	return context.WithValue(
		ctx,
		&userInfoKey,
		requestUserInfo{
			token:  token,
			userID: userID,
		},
	)
}

func GetUserID(ctx context.Context) (int, error) {
	userInfo, err := getUserInfo(ctx)
	if err != nil {
		return 0, err
	}
	return userInfo.userID, nil
}

func GetAccessToken(ctx context.Context) (string, error) {
	userInfo, err := getUserInfo(ctx)
	if err != nil {
		return "", err
	}
	return userInfo.token, nil
}

func getUserInfo(ctx context.Context) (requestUserInfo, error) {
	v := ctx.Value(&userInfoKey)
	userInfo, ok := v.(requestUserInfo)
	if !ok {
		return requestUserInfo{}, errors.New("unset userID")
	}
	return userInfo, nil
}
