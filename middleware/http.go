package middleware

import (
	"context"
	stdHttp "net/http"
	"os"

	"github.com/airunny/wiki-go-tools/icontext"
	"github.com/airunny/wiki-go-tools/iheader"
	"github.com/airunny/wiki-go-tools/reqid"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"go.opentelemetry.io/otel/trace"
)

var (
	allowOrigins = []string{"*"}
	allowHeaders = []string{"X-Symbol", "Authorization", "Content-Type", "X-User-Id"}
	allowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
)

func CORS() http.FilterFunc {
	return handlers.CORS(
		handlers.AllowedOrigins(allowOrigins),
		handlers.AllowedHeaders(allowHeaders),
		handlers.AllowedMethods(allowMethods),
		handlers.OptionStatusCode(204),
	)
}

func TraceIdAndRequestIdWithHeader(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var (
			traceID   string
			requestId string
		)

		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			traceID = span.TraceID().String()
		}

		if traceID != "" {
			tr.ReplyHeader().Set(iheader.TraceIdHeaderKey, traceID)
		}

		requestId = iheader.GetRequestId(tr.RequestHeader())
		if requestId == "" {
			requestId = reqid.GenRequestID()
		}

		tr.ReplyHeader().Set(iheader.RequestIdKey, requestId)
		ctx = icontext.WithRequestId(ctx, requestId)

		return handler(ctx, req)
	}
}

func TryParseHeader(opts ...Option) middleware.Middleware {
	o := Options{}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			header := tr.RequestHeader()
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
			// wsc
			wscValue := iheader.GetRouteWSC(header)
			ctx = icontext.WithWSC(ctx, wscValue)
			if wscValue != "" {
				tr.ReplyHeader().Set("route-wsc-env", os.Getenv("APOLLO_CLUSTER"))
			}
			return handler(ctx, req)
		}
	}
}

func ResponseEncoder(w http.ResponseWriter, r *stdHttp.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		stdHttp.Redirect(w, r, url, code)
		return nil
	}

	if res, ok := v.(TextPlainReply); ok {
		w.Header().Set("Content-Type", iheader.ResponseContentTextType)
		_, err := w.Write([]byte(res.StringReply()))
		if err != nil {
			w.WriteHeader(stdHttp.StatusInternalServerError)
		}
		return nil
	}

	WriteResponse(w, r, ResponseWithData(v))
	return nil
}

func ErrorEncoder(w http.ResponseWriter, r *stdHttp.Request, err error) {
	WriteResponse(w, r, ResponseWithError(errors.FromError(err)))
}

func WriteResponse(w http.ResponseWriter, _ *stdHttp.Request, body interface{}) {
	codec := encoding.GetCodec(json.Name)
	data, err := codec.Marshal(body)
	if err != nil {
		w.WriteHeader(stdHttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", iheader.ResponseContentJsonType)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(stdHttp.StatusInternalServerError)
	}
}
