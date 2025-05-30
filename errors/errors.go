package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrBadRequest          = errors.New(400, "INVALID_ARGS", "INVALID ARGS")
	ErrRefreshTokenExpired = errors.New(400, "REFRESH_TOKEN_EXPIRED", "not auth")
	ErrAccessTokenExpired  = errors.New(400, "ACCESS_TOKEN_EXPIRED", "token expired")
	ErrLogin               = errors.New(401, "UNAUTHORIZED", "not auth")
	ErrUserOperation       = errors.New(402, "USER_OPERATION", "try again later")
	ErrResourceNotFound    = errors.New(404, "RESOURCE_NOT_FOUND", "resource not found")
	ErrInternalServer      = errors.New(500, "INTERNAL_SERVER_ERROR", "internal server err")
)

func WithReason(e *errors.Error, in string) *errors.Error {
	return errors.New(int(e.Code), in, e.Message)
}

func WithMessage(e *errors.Error, msg string) *errors.Error {
	return errors.New(int(e.Code), e.Reason, msg)
}
