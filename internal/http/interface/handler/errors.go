package handler

import "examples/internal/http/interface/infra"

var (
	ErrUnexpected   = newErrorResponse("エラー", "不明なエラー")
	ErrPathNotExist = newErrorResponse("アクセスエラー", "アクセスパスが存在しません")
	ErrValidParam   = newErrorResponse("入力値エラー", "入力された値が不正です")
	ErrFailedSignin = newErrorResponse("サインイン失敗", "ログインIDまたはパスワードが間違っています")
)

func newErrorResponse(title, body string) *infra.ErrorResponse {
	return &infra.ErrorResponse{
		Title: title,
		Body:  body,
	}
}
