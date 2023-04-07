package handler

import (
	"context"
	"examples/pkg/adapter/framework/chi/infra"
	"examples/pkg/adapter/handler"
	cInfra "examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"examples/pkg/logic"
	"examples/pkg/logic/iodata"
	"examples/pkg/util"
	"fmt"
	"net/http"
	"strconv"
)

type UserHandler interface {
	Signup(ctx context.Context, c infra.ChiContext) *cInfra.HandleError
	GetAll(ctx context.Context, c infra.ChiContext) *cInfra.HandleError
	GetByID(ctx context.Context, c infra.ChiContext) *cInfra.HandleError
	ModifyAuthority(ctx context.Context, c infra.ChiContext) *cInfra.HandleError
	ModifyName(ctx context.Context, c infra.ChiContext) *cInfra.HandleError
}

type userHandler struct {
	userLogic logic.UserLogic
}

func NewUserHandler(userLogic logic.UserLogic) UserHandler {
	return &userHandler{
		userLogic: userLogic,
	}
}

func (h *userHandler) Signup(ctx context.Context, c infra.ChiContext) *cInfra.HandleError {
	var input *iodata.SignupInput
	if err := c.Decode(&input); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	if err := h.userLogic.Signup(ctx, input); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) GetAll(ctx context.Context, c infra.ChiContext) *cInfra.HandleError {
	users, err := h.userLogic.GetAll(ctx)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			r = handler.NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	// ここの変換処理は Presenter が本来担当する
	var output []iodata.UserOutput
	for _, user := range users {
		output = append(output, iodata.UserOutput{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		})
	}
	if err := c.WriteJSON(http.StatusOK, output); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) GetByID(ctx context.Context, c infra.ChiContext) *cInfra.HandleError {
	userID, _ := strconv.Atoi(c.URLParam("user_id"))
	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			r = handler.NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	if err := c.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) ModifyAuthority(ctx context.Context, c infra.ChiContext) *cInfra.HandleError {
	// リクエスト者の権限を確認
	const requiredAuthority = 99
	if ok, err := h.userLogic.Authorization(ctx, requiredAuthority); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	} else if !ok {
		err = errors.Errorf(code.CodeUnauthorized, "lack of authority: %d", requiredAuthority)
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	var input *iodata.ModifyAuthorityInput
	if err := c.Decode(&input); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}
	if err := input.Validate(); err != nil {
		return &cInfra.HandleError{}
	}

	userID, _ := strconv.Atoi(c.URLParam("user_id"))
	if err := h.userLogic.ModifyAuthority(ctx, userID, input.Authority); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	// 更新後データ 再取得
	// この時点で更新は完了しているので、DBアクセスで何かしらのエラーが発生して
	// エラーが返却されていても、ログ出力のみに留め正常終了扱いでレスポンスを返却する
	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		logger.L.Warn(fmt.Sprint(err))
	}

	if err := c.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) ModifyName(ctx context.Context, c infra.ChiContext) *cInfra.HandleError {
	userID, err := util.GetUserID(ctx)
	if err != nil {
		return &cInfra.HandleError{}
	}

	var input *iodata.ModifyNameInput
	if err := c.Decode(&input); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}
	if err := input.Validate(); err != nil {
		return &cInfra.HandleError{}
	}

	if err := h.userLogic.ModifyName(ctx, userID, input.Name); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	// 更新後データ 再取得
	// この時点で更新は完了しているので、DBアクセスで何かしらのエラーが発生して
	// エラーが返却されていても、ログ出力のみに留め正常終了扱いでレスポンスを返却する
	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		logger.L.Warn(fmt.Sprint(err))
	}

	if err := c.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}
