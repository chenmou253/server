package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type OpenAPIHeaders struct {
	APIKey    string `json:"apiKey"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

func BuildOpenAPISignature(secret string, apiKey string, timestamp string, nonce string, method string, path string, body []byte) string {
	bodyHash := sha256.Sum256(body)
	raw := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
		apiKey,
		timestamp,
		nonce,
		strings.ToUpper(method),
		path,
		hex.EncodeToString(bodyHash[:]),
		secret,
	)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func NewOpenAPIHeaders(apiKey string, secret string, method string, path string, body []byte) (*OpenAPIHeaders, error) {
	nonce, err := generateNonce(16)
	if err != nil {
		return nil, err
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	return &OpenAPIHeaders{
		APIKey:    apiKey,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: BuildOpenAPISignature(secret, apiKey, timestamp, nonce, method, path, body),
	}, nil
}

func generateNonce(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
