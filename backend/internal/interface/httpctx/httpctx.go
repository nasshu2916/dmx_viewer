package httpctx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	keyRequestID contextKey = "request_id"
	keyRealIP    contextKey = "real_ip"
)

// WithRequestID はコンテキストへ Request-ID を格納する
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID はコンテキストから Request-ID を取得する
func RequestID(ctx context.Context) string {
	if v := ctx.Value(keyRequestID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// WithRealIP はコンテキストへ Real-IP を格納する
func WithRealIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyRealIP, ip)
}

// RealIP はコンテキストから Real-IP を取得する
func RealIP(ctx context.Context) string {
	if v := ctx.Value(keyRealIP); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RequestIDFromHeaderOrNew はヘッダに Request-ID があれば利用し、なければ新規発行する
func RequestIDFromHeaderOrNew(r *http.Request) string {
	if id := r.Header.Get("X-Request-Id"); id != "" {
		return id
	}
	if id := r.Header.Get("Request-Id"); id != "" {
		return id
	}
	// ランダム16バイトをhex化(32桁)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	// フォールバック(時刻ベース)
	return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
}

// ResolveRealIP は X-Forwarded-For / X-Real-IP / RemoteAddr からクライアントIPを解決する
func ResolveRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
