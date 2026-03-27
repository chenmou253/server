# Open API

## Overview

Public APIs are exposed under `/api/v1/open/*`.

These APIs:

- do not use JWT
- use `X-API-Key + X-Timestamp + X-Nonce + X-Signature`
- are intended for third-party integration

Base URL example:

```text
http://127.0.0.1:8080
```

## Auth Headers

Every request must include:

- `X-API-Key`
- `X-Timestamp`
- `X-Nonce`
- `X-Signature`
- `Content-Type: application/json`

## Env Config

Configured in [server/.env](/Users/jiechen/go-admin/server/.env):

```env
OPEN_API_KEY=replace-with-open-api-key
OPEN_API_SECRET=replace-with-open-api-secret
OPEN_API_TIME_SKEW_SEC=300
```

## Signature Rule

Signature source string:

```text
apiKey
timestamp
nonce
METHOD
/path
sha256(body)
secret
```

Then:

1. SHA256 the request body, lowercase hex
2. Put that hash into the source string
3. SHA256 the whole source string
4. Put the lowercase hex result into `X-Signature`

Notes:

- use uppercase `METHOD`
- sign path only, no domain
- query string is not included in signature
- for empty body, use SHA256 of empty string

Empty body hash:

```text
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Public Device APIs

### 1. Register Device

`POST /api/v1/open/devices`

Request body:

```json
{
  "deviceNo": "DEV001",
  "merchantId": "M1001",
  "status": 0,
  "ip": "192.168.1.10",
  "createT": 1710000000
}
```

Fields map directly to table columns:

- `deviceNo` -> `device_no`
- `merchantId` -> `merchant_id`
- `status` -> `status`
- `ip` -> `ip`
- `createT` -> `create_t`

Notes:

- `id` is auto-increment, do not send it
- if `createT` is omitted or `0`, server will use current Unix timestamp
- `status` only allows `0` or `1`

Response example:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "deviceNo": "DEV001",
    "merchantId": "M1001",
    "status": 0,
    "ip": "192.168.1.10",
    "createT": 1710000000
  }
}
```

### 2. Query Device List

`GET /api/v1/open/devices`

Query params:

- `page`
- `pageSize`
- `deviceNo`
- `merchantId`
- `status`

### 3. Query Single Device

`GET /api/v1/open/devices/{deviceNo}`

### 4. Update Device Status

`POST /api/v1/open/devices/{deviceNo}/status`

Request body:

```json
{
  "status": 1
}
```

## cURL Example

### Register Device

```bash
API_KEY="your-api-key"
SECRET="your-api-secret"
TIMESTAMP=$(date +%s)
NONCE="reg-001"
BODY='{"deviceNo":"DEV001","merchantId":"M1001","status":0,"ip":"192.168.1.10","createT":1710000000}'
BODY_SHA256=$(printf '%s' "$BODY" | shasum -a 256 | awk '{print $1}')
PATH_ONLY="/api/v1/open/devices"

SIGN_RAW="${API_KEY}
${TIMESTAMP}
${NONCE}
POST
${PATH_ONLY}
${BODY_SHA256}
${SECRET}"

SIGNATURE=$(printf '%s' "$SIGN_RAW" | shasum -a 256 | awk '{print $1}')

curl -X POST "http://127.0.0.1:8080${PATH_ONLY}" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: ${API_KEY}" \
  -H "X-Timestamp: ${TIMESTAMP}" \
  -H "X-Nonce: ${NONCE}" \
  -H "X-Signature: ${SIGNATURE}" \
  -d "$BODY"
```

## Go Example

```go
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func buildSignature(apiKey, secret, timestamp, nonce, method, path string, body []byte) string {
	raw := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
		apiKey,
		timestamp,
		nonce,
		method,
		path,
		sha256Hex(body),
		secret,
	)
	return sha256Hex([]byte(raw))
}

func main() {
	apiKey := "your-api-key"
	secret := "your-api-secret"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := "go-reg-001"
	path := "/api/v1/open/devices"
	body := []byte(`{"deviceNo":"DEV001","merchantId":"M1001","status":0,"ip":"192.168.1.10","createT":1710000000}`)
	signature := buildSignature(apiKey, secret, timestamp, nonce, "POST", path, body)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080"+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}
```

## Related Code

- [server/internal/middleware/openapi.go](/Users/jiechen/go-admin/server/internal/middleware/openapi.go)
- [server/internal/handler/public.go](/Users/jiechen/go-admin/server/internal/handler/public.go)
- [server/internal/service/public.go](/Users/jiechen/go-admin/server/internal/service/public.go)
- [server/internal/router/router.go](/Users/jiechen/go-admin/server/internal/router/router.go)
