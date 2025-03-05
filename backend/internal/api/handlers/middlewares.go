package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"notes/internal/configs"
	"notes/pkg/date"
	"notes/pkg/response"
	"notes/pkg/validations"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/tomasen/realip"
)

func (h *Handlers) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					h.HttpErrs.ServerError(w, r, fmt.Errorf("%s", err), APIErrKey)
				}
			}()
			next.ServeHTTP(w, r)
		})
}

func (h *Handlers) LogAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsWriter(w)
		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)
		next.ServeHTTP(mw, r)
		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)
		h.Logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (h *Handlers) AddHeadersWithCSP(next http.Handler) http.Handler {
	conf := configs.New()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			if handlePreflight(w, origin, conf) {
				return
			}
		}

		if !handleActualRequest(w, origin, conf) {
			return
		}

		addSecurityHeaders(w)
		addContentSecurityPolicy(w, r, origin)
		next.ServeHTTP(w, r)
	})
}

func handlePreflight(w http.ResponseWriter, origin string, conf *configs.Config) bool {
	log.Printf("Handling preflight for origin: %s", origin)

	if conf.ENV == "production" && !isValidOrigin(origin, conf) {
		log.Printf("Blocked invalid preflight origin: %s", origin)
		http.Error(w, "Invalid Origin: "+origin, http.StatusForbidden)
		return true
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusNoContent)
	return true
}

func handleActualRequest(w http.ResponseWriter, origin string, conf *configs.Config) bool {
	log.Printf("Handling actual request from origin: %s", origin)

	if !isValidOrigin(origin, conf) && conf.ENV == "production" {
		log.Printf("Blocked invalid origin: %s", origin)
		http.Error(w, "Invalid Origin: "+origin, http.StatusForbidden)
		return false
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	return true
}

func isValidOrigin(origin string, conf *configs.Config) bool {
	if origin == "" {
		return true
	}

	allowedOrigins := strings.Split(strings.ReplaceAll(conf.ALLOWED_ORIGINS, " ", ""), ",")
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

func addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
}

func addContentSecurityPolicy(w http.ResponseWriter, r *http.Request, origin string) {
	if strings.HasPrefix(r.URL.Path, "/swagger/") {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
				"img-src data: https://cdn.ngrok.com; "+
				"connect-src 'self' "+origin)
	} else {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"connect-src 'self' "+origin+"; "+
				"img-src 'self' data:; "+
				"script-src 'self' 'unsafe-inline'")
	}
}

// DEBUG END
func (h *Handlers) WithTimeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type contextKey string

const userIDKey contextKey = "userID"

type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (h *Handlers) PROTECT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(configs.GetString("JWT_NAME", ""))
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			h.HttpErrs.CheckErrType(w, r, validations.ErrJWT)
			return
		}

		tokenStr := cookie.Value

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.GetString("JWT_STRING", "")), nil
		})
		if err != nil || !token.Valid {
			h.HttpErrs.CheckErrType(w, r, validations.ErrJWT)
			return
		}

		if claims.ExpiresAt != nil && date.ArgentinaTimeNow().After(claims.ExpiresAt.Time) {
			h.HttpErrs.CheckErrType(w, r, validations.ErrTokenExpired)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (*uint, error) {
	if userID, ok := ctx.Value(userIDKey).(uint); ok {
		return &userID, nil
	}
	return nil, validations.ErrUnauthorized
}

func (h *Handlers) LimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := rl.AddVisitor(ip)

			if !limiter.Allow() {
				h.HttpErrs.CheckErrType(w, r, validations.ErrRateLimitExcess)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handlers) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := response.NewMetricsWriter(w)

		next.ServeHTTP(mw, r)

		if r.URL.Path == "/prometheus-metrics" {
			return
		}

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(mw.StatusCode)

		// Get route pattern from mux
		route := "unknown"
		if currentRoute := mux.CurrentRoute(r); currentRoute != nil {
			if path, err := currentRoute.GetPathTemplate(); err == nil {
				route = path
			}
		}

		h.HttpDuration.WithLabelValues(r.Method, route).Observe(duration)
		h.HttpRequestsTotal.WithLabelValues(r.Method, route, status).Inc()
	})
}
