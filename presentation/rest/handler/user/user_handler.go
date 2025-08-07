package handler

import (
	"github.com/gin-gonic/gin"
	userPort "github.com/pdh9523/go-hexarch/domains/user/application/port/in"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/request"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"github.com/pdh9523/go-hexarch/presentation/rest/mapper"
)

type UserHandler struct {
	userUseCase userPort.UserUseCase
	responder   *response.Responder
}

func NewUserHandler(userUseCase userPort.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		responder:   response.NewResponder(),
	}
}

func (h *UserHandler) GetMyInfo(c *gin.Context) {
	result, err := h.userUseCase.GetMyInfo(c.Request.Context())
	if err != nil {
		h.responder.Error(c, err)
		return
	}
	h.responder.Success(c, result)
}

func (h *UserHandler) CheckNicknameAvailability(c *gin.Context) {
	var req request.CheckNicknameAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	query := mapper.ToCheckNicknameQuery(req)
	err := h.userUseCase.CheckNicknameAvailability(c.Request.Context(), query)
	if err != nil {
		h.responder.Error(c, err)
		return
	}
	h.responder.Success(c, nil)
}

func (h *UserHandler) CheckUsernameAvailability(c *gin.Context) {
	var req request.CheckUsernameAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	query := mapper.ToCheckUsernameQuery(req)
	err := h.userUseCase.CheckUsernameAvailability(c.Request.Context(), query)
	if err != nil {
		h.responder.Error(c, err)
		return
	}
	h.responder.Success(c, nil)
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	command := mapper.ToSignUpCommand(req)
	result, err := h.userUseCase.SignUp(c.Request.Context(), command)
	if err != nil {
		h.responder.Error(c, err)
		return
	}
	res := mapper.ToTokenResponse(result)
	h.responder.Success(c, res)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	command := mapper.ToSignInCommand(req)
	result, err := h.userUseCase.SignIn(c.Request.Context(), command)
	if err != nil {
		h.responder.Error(c, err)
		return
	}
	res := mapper.ToTokenResponse(result)
	h.responder.Success(c, res)
}

func (h *UserHandler) SignOut(c *gin.Context) {
	if err := h.userUseCase.SignOut(c.Request.Context()); err != nil {
		h.responder.Error(c, err)
		return
	}
	h.responder.Success(c, nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	command := mapper.ToChangePasswordCommand(req)
	err := h.userUseCase.ChangePassword(c.Request.Context(), command)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, nil)
}

func (h *UserHandler) ChangeNickname(c *gin.Context) {
	var req request.ChangeNicknameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responder.Error(c, err)
		return
	}

	command := mapper.ToChangeNicknameCommand(req)
	err := h.userUseCase.ChangeNickname(c.Request.Context(), command)
	if err != nil {
		h.responder.Error(c, err)
		return
	}

	h.responder.Success(c, nil)
}
