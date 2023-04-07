package handler

import "examples/pkg/adapter/infra"

var (
	ErrUnexpected       = NewHTTPError("エラー", "不明なエラー")
	ErrPathNotExist     = NewHTTPError("アクセスエラー", "アクセスパスが存在しません")
	ErrValidParam       = NewHTTPError("入力値エラー", "入力された値が不正です")
	ErrFailedSignin     = NewHTTPError("サインイン失敗", "ログインIDまたはパスワードが間違っています")
	HTTPErrUnauthorized = NewHTTPError("セッションエラー", "セッションが切れました。\n再度サインインしてください。")
)

func NewHTTPError(title, body string) *infra.HTTPError {
	return &infra.HTTPError{
		Title: title,
		Body:  body,
	}
}
