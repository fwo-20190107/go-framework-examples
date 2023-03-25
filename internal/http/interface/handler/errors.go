package handler

import "examples/internal/http/interface/infra"

var (
	ErrUnexpected   = NewHTTPError("エラー", "不明なエラー")
	ErrPathNotExist = NewHTTPError("アクセスエラー", "アクセスパスが存在しません")
	ErrValidParam   = NewHTTPError("入力値エラー", "入力された値が不正です")
	ErrFailedSignin = NewHTTPError("サインイン失敗", "ログインIDまたはパスワードが間違っています")
)

func NewHTTPError(title, body string) *infra.HTTPError {
	return &infra.HTTPError{
		Title: title,
		Body:  body,
	}
}
