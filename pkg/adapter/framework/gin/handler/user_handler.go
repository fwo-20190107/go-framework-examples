package handler

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/adapter/infra"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"examples/pkg/logic"
	"examples/pkg/logic/iodata"
	"examples/pkg/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Signup(c *gin.Context) *infra.HandleError
	GetAll(c *gin.Context) *infra.HandleError
	GetByID(c *gin.Context) *infra.HandleError
	ModifyAuthority(c *gin.Context) *infra.HandleError
	ModifyName(c *gin.Context) *infra.HandleError
}

type userHandler struct {
	userLogic logic.UserLogic
}

func NewUserHandler(userLogic logic.UserLogic) UserHandler {
	return &userHandler{
		userLogic: userLogic,
	}
}

func (h *userHandler) Signup(c *gin.Context) *infra.HandleError {
	var input *iodata.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	if err := input.Validate(); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}

	if err := h.userLogic.Signup(c, input); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}
	return nil
}

func (h *userHandler) GetAll(c *gin.Context) *infra.HandleError {
	users, err := h.userLogic.GetAll(c)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			r = handler.NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &infra.HandleError{HTTPError: r, Error: err}
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

	c.JSON(http.StatusOK, output)
	return nil
}

func (h *userHandler) GetByID(c *gin.Context) *infra.HandleError {
	uidp := c.Param("user_id")
	userID, err := strconv.Atoi(uidp)
	if err != nil {

	}

	user, err := h.userLogic.GetByID(c, userID)
	if err != nil {
		r := handler.ErrUnexpected
		switch {
		case errors.Is(err, code.CodeNotFound):
			r = handler.NewHTTPError("エラー", "ユーザーデータなし")
		}
		return &infra.HandleError{HTTPError: r, Error: err}
	}

	c.JSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	})
	return nil
}

func (h *userHandler) ModifyAuthority(c *gin.Context) *infra.HandleError {
	// リクエスト者の権限を確認
	const requiredAuthority = 99
	if ok, err := h.userLogic.Authorization(c, requiredAuthority); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	} else if !ok {
		err = errors.Errorf(code.CodeUnauthorized, "lack of authority: %d", requiredAuthority)
		return &infra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	uidp := c.Param("user_id")
	userID, err := strconv.Atoi(uidp)
	if err != nil {

	}

	var input *iodata.ModifyAuthorityInput
	if err := c.BindJSON(&input); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}
	if err := input.Validate(); err != nil {
		return &infra.HandleError{}
	}

	if err := h.userLogic.ModifyAuthority(c, userID, input.Authority); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	// 更新後データ 再取得
	// この時点で更新は完了しているので、DBアクセスで何かしらのエラーが発生して
	// エラーが返却されていても、ログ出力のみに留め正常終了扱いでレスポンスを返却する
	user, err := h.userLogic.GetByID(c, userID)
	if err != nil {
		logger.L.Warn(c, fmt.Sprint(err))
	}

	c.JSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	})
	return nil
}

func (h *userHandler) ModifyName(c *gin.Context) *infra.HandleError {
	userID, err := util.GetUserID(c)
	if err != nil {
		return &infra.HandleError{}
	}

	var input *iodata.ModifyNameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrValidParam, Error: err}
	}
	if err := input.Validate(); err != nil {
		return &infra.HandleError{}
	}

	if err := h.userLogic.ModifyName(c, userID, input.Name); err != nil {
		return &infra.HandleError{HTTPError: handler.ErrUnexpected, Error: err}
	}

	// 更新後データ 再取得
	// この時点で更新は完了しているので、DBアクセスで何かしらのエラーが発生して
	// エラーが返却されていても、ログ出力のみに留め正常終了扱いでレスポンスを返却する
	user, err := h.userLogic.GetByID(c, userID)
	if err != nil {
		logger.L.Warn(c, fmt.Sprint(err))
	}

	c.JSON(http.StatusOK, iodata.UserOutput{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	})
	return nil
}
