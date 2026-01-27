package middlewares

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestIDMiddleware adds a unique request ID to each request.
// The ID is available in the request context and response headers.
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check if request already has an ID (from upstream proxy)
			requestID := c.Request().Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// set request ID in context and response header
			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}

// RequestLoggerMiddleware logs HTTP requests with zerolog.
// It includes request ID, method, path, status, latency, and other useful info.
func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// process request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// gather request info
			req := c.Request()
			res := c.Response()
			latency := time.Since(start)

			// get request ID from context
			requestID, _ := c.Get(RequestIDKey).(string)

			// log the request
			event := log.Info()
			if res.Status >= 500 {
				event = log.Error()
			} else if res.Status >= 400 {
				event = log.Warn()
			}

			event.
				Str("request_id", requestID).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("latency", latency).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Int64("bytes_in", req.ContentLength).
				Int64("bytes_out", res.Size).
				Msg("request")

			return nil
		}
	}
}

// GetRequestID retrieves the request ID from the Echo context
func GetRequestID(c echo.Context) string {
	if id, ok := c.Get(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
