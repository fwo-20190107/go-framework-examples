package util

import (
	"context"
	"errors"
	"examples/internal/gin/entity"
)

type userInfo struct {
	token  string
	userID int
}

type ctxkey struct{}

var userInfoKey = ctxkey{}

func newUserInfo(token string, userID int) userInfo {
	return userInfo{
		token:  token,
		userID: userID,
	}
}

func SetUserInfo(ctx context.Context, token string, userID int) context.Context {
	return context.WithValue(
		ctx,
		&entity.UserInfoKey,
		newUserInfo(token, userID),
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

func getUserInfo(ctx context.Context) (userInfo, error) {
	info, ok := ctx.Value(&userInfoKey).(userInfo)
	if !ok {
		return userInfo{}, errors.New("unset userID")
	}
	return info, nil
}
