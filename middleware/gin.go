package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/airunny/wiki-go-tools/icontext"
	"github.com/airunny/wiki-go-tools/iheader"
	"github.com/airunny/wiki-go-tools/reqid"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func TraceIdAndRequestIdWithHeaderForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx       = c.Request.Context()
			traceID   string
			requestId string
		)

		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			traceID = span.TraceID().String()
		}

		if traceID != "" {
			c.Writer.Header().Set(iheader.TraceIdHeaderKey, traceID)
		}

		requestId = iheader.GetRequestId(headerCarrier(c.Request.Header))
		if requestId == "" {
			requestId = reqid.GenRequestID()
		}

		c.Writer.Header().Set(iheader.RequestIdKey, requestId)
		ctx = icontext.WithRequestId(ctx, requestId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TryParseHeaderForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx    = c.Request.Context()
			header = headerCarrier(c.Request.Header)
		)

		// 客户端ip
		ctx = icontext.WithClientIP(ctx, iheader.GetClientIp(header))
		// 用户ID
		ctx = icontext.WithUserId(ctx, iheader.GetUserId(header))
		// basic data
		ctx = icontext.WithBasicData(ctx, iheader.GetBasicData(header))
		// 城市码
		countryCode := iheader.GetCountryCode(header)
		ctx = icontext.WithCountryCode(ctx, countryCode)
		// 语言code
		ctx = icontext.WithLanguageCode(ctx, iheader.GetLanguageCode(header))
		// 偏好语言
		ctx = icontext.WithPreferredLanguageCode(ctx, iheader.GetPreferredLanguageCode(header))
		// wiki data center Request-Id
		ctx = icontext.WithWikiDataCenterRequestId(ctx, iheader.GetWikiDataCenterRequestId(header))
		// scene code
		ctx = icontext.WithSceneCode(ctx, iheader.GetSceneCode(header))
		// x-pwa
		ctx = icontext.WithXPWA(ctx, iheader.GetXPwa(header))
		// wsc
		wscValue := iheader.GetRouteWSC(header)
		ctx = icontext.WithWSC(ctx, wscValue)
		if wscValue != "" {
			c.Writer.Header().Set("route-wsc-env", os.Getenv("APOLLO_CLUSTER"))
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

func (hc headerCarrier) Add(key string, value string) {
	http.Header(hc).Add(key, value)
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

var ginResultKey = "res"

func JSONData(c *gin.Context, data interface{}) {
	c.Set(ginResultKey, data)
	c.Next()
}

func JSONError(c *gin.Context, err error) {
	e := errors.FromError(err)
	c.JSON(http.StatusOK, BizResponse{
		Code:    e.Code,
		Message: e.Message,
		Reason:  e.Reason,
		Time:    time.Now().Unix(),
	})
	c.Abort()
}

func JSONResponse(c *gin.Context) {
	ret, _ := c.Get(ginResultKey)
	c.JSON(http.StatusOK, ResponseWithData(ret))
}
