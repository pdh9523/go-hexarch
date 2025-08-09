package mapper

import (
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/command"
	"github.com/pdh9523/go-hexarch/domains/user/application/port/in/query"
	"github.com/pdh9523/go-hexarch/domains/user/domain/result"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/request"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
)

func ToCheckNicknameQuery(req request.CheckNicknameAvailabilityRequest) query.CheckNicknameQuery {
	return query.CheckNicknameQuery{
		Nickname: req.Nickname,
	}
}

func ToCheckUsernameQuery(req request.CheckUsernameAvailabilityRequest) query.CheckUsernameQuery {
	return query.CheckUsernameQuery{
		Username: req.Username,
	}
}

func ToSignUpCommand(req request.SignUpRequest) command.SignUpCommand {
	return command.SignUpCommand{
		Nickname: req.Nickname,
		Username: req.Username,
		Password: req.Password,
	}
}

func ToSignInCommand(req request.SignInRequest) command.SignInCommand {
	return command.SignInCommand{
		Username: req.Username,
		Password: req.Password,
	}
}

func ToChangePasswordCommand(req request.ChangePasswordRequest) command.ChangePasswordCommand {
	return command.ChangePasswordCommand{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}
}

func ToTokenResponse(result *result.TokenResult) response.TokenResponse {
	return response.TokenResponse{
		AccessToken:  result.Token.AccessToken,
		RefreshToken: result.Token.RefreshToken,
	}
}

func ToChangeNicknameCommand(req request.ChangeNicknameRequest) command.ChangeNicknameCommand {
	return command.ChangeNicknameCommand{
		Nickname: req.Nickname,
	}
}

func ToUserInfoResponse(result *result.UserInfoResult) response.UserInfoResponse {
	return response.UserInfoResponse{
		Username: result.Username,
		Nickname: result.Nickname,
		Role:     result.Role,
	}
}

func ToRefreshTokenCommand(req request.RefreshTokenRequest) command.RefreshTokenCommand {
	return command.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}
}
