package mapper

import (
	"errors"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	userError "github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/pdh9523/go-hexarch/infrastructure/gorm/model"
	commonError "github.com/pdh9523/go-hexarch/shared/common/errorCode"
	"gorm.io/gorm"
)

func ToDomainUser(userModel *model.User) (*domain.User, error) {
	role := stringToRole(userModel.Role)

	return &domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Nickname: userModel.Nickname,
		Password: userModel.Password,
		Role:     role,
	}, nil
}

func ToModelUser(userDomain *domain.User) *model.User {
	role := string(userDomain.Role)

	user := &model.User{
		Username: userDomain.Username,
		Nickname: userDomain.Nickname,
		Password: userDomain.Password,
		Role:     role,
	}

	if userDomain.ID != "" {
		user.ID = userDomain.ID
	}

	return user
}

func stringToRole(role string) domain.Role {
	var result domain.Role
	switch role {
	case "ADMIN":
		result = domain.ROLE_ADMIN
	case "USER":
		result = domain.ROLE_USER
	default:
		result = domain.ROLE_USER
	}
	return result
}

func ToDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return userError.ErrUserNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return userError.ErrUserAlreadyExists
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return commonError.ErrInternalServerError
	case errors.Is(err, gorm.ErrNotImplemented):
		return commonError.ErrInternalServerError
	case errors.Is(err, gorm.ErrMissingWhereClause):
		return commonError.ErrInvalidRequest
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return commonError.ErrInternalServerError
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return commonError.ErrInvalidRequest
	case errors.Is(err, gorm.ErrModelValueRequired):
		return commonError.ErrInvalidRequest
	case errors.Is(err, gorm.ErrInvalidData):
		return commonError.ErrInvalidRequest
	default:
		// 알 수 없는 DB 에러는 내부 서버 에러로 처리
		return commonError.ErrInternalServerError
	}
}
