package middleware

import (
	"bytes"
	"crypto/subtle"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/config"
	"go-admin/server/internal/contextx"
	"go-admin/server/internal/utils"
	"go-admin/server/pkg/response"
)

const (
	openAPIKeyHeader       = "X-API-Key"
	openAPITimestampHeader = "X-Timestamp"
	openAPINonceHeader     = "X-Nonce"
	openAPISignatureHeader = "X-Signature"
)

func OpenAPIAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := strings.TrimSpace(c.GetHeader(openAPIKeyHeader))
		timestamp := strings.TrimSpace(c.GetHeader(openAPITimestampHeader))
		nonce := strings.TrimSpace(c.GetHeader(openAPINonceHeader))
		signature := strings.TrimSpace(c.GetHeader(openAPISignatureHeader))

		if apiKey == "" || timestamp == "" || nonce == "" || signature == "" {
			response.Error(c, http.StatusUnauthorized, "missing open api auth headers")
			c.Abort()
			return
		}
		if apiKey != cfg.OpenAPIKey {
			response.Error(c, http.StatusUnauthorized, "invalid api key")
			c.Abort()
			return
		}

		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid timestamp")
			c.Abort()
			return
		}
		now := time.Now().Unix()
		if absInt64(now-ts) > int64(cfg.OpenAPITimeSkewSec) {
			response.Error(c, http.StatusUnauthorized, "timestamp expired")
			c.Abort()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "read body failed")
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		expected := utils.BuildOpenAPISignature(cfg.OpenAPISecret, apiKey, timestamp, nonce, c.Request.Method, c.Request.URL.Path, bodyBytes)
		if subtle.ConstantTimeCompare([]byte(strings.ToLower(signature)), []byte(expected)) != 1 {
			response.Error(c, http.StatusUnauthorized, "invalid signature")
			c.Abort()
			return
		}

		c.Set(contextx.OpenAPIKey, apiKey)
		c.Next()
	}
}

func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
