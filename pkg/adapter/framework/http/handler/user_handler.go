package handler

import (
	"context"
	"examples/pkg/adapter/framework/http/infra"
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
	Signup(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError
	HandleRoot(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError
}

type userHandler struct {
	userLogic logic.UserLogic
}

func NewUserHandler(userLogic logic.UserLogic) UserHandler {
	return &userHandler{
		userLogic: userLogic,
	}
}

// signup godoc
//	@Summary		Sign up of the application
//	@Description	Register account information and create user
//	@tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input	body	iodata.SignupInput	true	"foo"
//	@Success		200
//	@Failure		400	{object}	infra.HTTPError
//	@Failure		401	{object}	infra.HTTPError
//	@Failure		404	{object}	infra.HTTPError
//	@Failure		500	{object}	infra.HTTPError
//	@Router			/signup [post]
func (h *userHandler) Signup(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
	if httpCtx.Method() != http.MethodPost {
		return &cInfra.HandleError{HTTPError: handler.ErrPathNotExist}
	}

	var input *iodata.SignupInput
	if err := httpCtx.Decode(&input); err != nil {
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

// getAll godoc
//	@Summary		Get all users
//	@Description	get all users
//	@tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]iodata.UserOutput
//	@Failure		400	{object}	infra.HTTPError
//	@Failure		401	{object}	infra.HTTPError
//	@Failure		404	{object}	infra.HTTPError
//	@Failure		500	{object}	infra.HTTPError
//	@Security		Bearer
//	@Router			/user [get]
func (h *userHandler) getAll(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
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
	if err := httpCtx.WriteJSON(http.StatusOK, output); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

// getByID godoc
//	@Summary		Get user by userID
//	@Description	get user
//	@tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int	true	"foo"
//	@Success		200		{object}	iodata.UserOutput
//	@Failure		400		{object}	infra.HTTPError
//	@Failure		401		{object}	infra.HTTPError
//	@Failure		404		{object}	infra.HTTPError
//	@Failure		500		{object}	infra.HTTPError
//	@Security		Bearer
//	@Router			/user/{user_id} [get]
func (h *userHandler) getByID(ctx context.Context, httpCtx infra.HttpContext, userID int) *cInfra.HandleError {
	user, err := h.userLogic.GetByID(ctx, userID)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			r = handler.NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &cInfra.HandleError{HTTPError: r, Error: err}
	}

	if err := httpCtx.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

// modifyAuthority godoc
//	@Summary		Modify user authority
//	@Description	Accepts authority changes from only admin user
//	@tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int							false	"modify target userID"
//	@Param			input	body		iodata.ModifyAuthorityInput	true	"foo"
//	@Success		200		{object}	iodata.UserOutput
//	@Failure		400		{object}	infra.HTTPError
//	@Failure		401		{object}	infra.HTTPError
//	@Failure		404		{object}	infra.HTTPError
//	@Failure		500		{object}	infra.HTTPError
//	@Security		Bearer
//	@Router			/user/{user_id} [patch]
func (h *userHandler) modifyAuthority(ctx context.Context, httpCtx infra.HttpContext, userID int) *cInfra.HandleError {
	// リクエスト者の権限を確認
	const requiredAuthority = 99
	if ok, err := h.userLogic.Authorization(ctx, requiredAuthority); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	} else if !ok {
		err = errors.Errorf(code.CodeUnauthorized, "lack of authority: %d", requiredAuthority)
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	var input *iodata.ModifyAuthorityInput
	if err := httpCtx.Decode(&input); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}
	if err := input.Validate(); err != nil {
		return &cInfra.HandleError{}
	}

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

	if err := httpCtx.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

// modifyName godoc
//	@Summary		Modify user name
//	@Description	Accepts name changes from the person himself
//	@tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input	body		iodata.ModifyNameInput	true	"foo"
//	@Success		200		{object}	iodata.UserOutput
//	@Failure		400		{object}	infra.HTTPError
//	@Failure		401		{object}	infra.HTTPError
//	@Failure		404		{object}	infra.HTTPError
//	@Failure		500		{object}	infra.HTTPError
//	@Security		Bearer
//	@Router			/user [patch]
func (h *userHandler) modifyName(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
	userID, err := util.GetUserID(ctx)
	if err != nil {
		return &cInfra.HandleError{}
	}

	var input *iodata.ModifyNameInput
	if err := httpCtx.Decode(&input); err != nil {
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

	if err := httpCtx.WriteJSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}); err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) HandleRoot(ctx context.Context, httpCtx infra.HttpContext) *cInfra.HandleError {
	vars, err := httpCtx.Vars("/user", "user_id")
	if err != nil {
		return &cInfra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	var userID int
	uidp, ok := vars["user_id"]
	if ok {
		if userID, err = strconv.Atoi(uidp); err != nil {
			return &cInfra.HandleError{HTTPError: handler.ErrValidParam, Error: errors.Errorf(code.CodeBadRequest, "path is not number")}
		}
	}

	switch httpCtx.Method() {
	case http.MethodGet:
		if !ok {
			return h.getAll(ctx, httpCtx)
		}
		return h.getByID(ctx, httpCtx, userID)
	case http.MethodPatch:
		if !ok {
			return h.modifyName(ctx, httpCtx)
		}
		return h.modifyAuthority(ctx, httpCtx, userID)
	}
	return nil
}
