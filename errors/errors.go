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
	ErrOutOfStock          = errors.New(403, "OUT_OF_STOCK", "out of stock")
	ErrResourceNotFound    = errors.New(404, "RESOURCE_NOT_FOUND", "resource not found")
	ErrPriceChanged        = errors.New(405, "PRICE_CHANGED", "price changed")
	ErrSellOut             = errors.New(406, "SELL_OUT", "sell out")
	ErrGoodsOff            = errors.New(407, "GOODS_OFF", "goods off")
	ErrBuyLimit            = errors.New(408, "BUY_LIMIT", "buy limit")
	ErrNotSupportDeliver   = errors.New(409, "NOT_SUPPORT_DELIVER", "address does not support delivery")
	ErrInternalServer      = errors.New(500, "INTERNAL_SERVER_ERROR", "internal server err")
)

func WithReason(e *errors.Error, in string) *errors.Error {
	return errors.New(int(e.Code), in, e.Message)
}

func WithMessage(e *errors.Error, msg string) *errors.Error {
	return errors.New(int(e.Code), e.Reason, msg)
}
