package gapi

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func GrpcLogger(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	duration := time.Since(time.Now())
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Str("protocol", "grpc").
		Msg("received a grpc request")
	return handler(ctx, req)
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: writer,
			statusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, request)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}
		logger.
			Str("protocol", "http").
			Str("method", request.Method).
			Str("path", request.RequestURI).
			Int("status_code", rec.statusCode).
			Str("status_text", http.StatusText(rec.statusCode)).
			Dur("duration", duration).
			Msg("received an http request")
	})
}
