package masking

import (
	"crypto/rand"
	"encoding/base64"
)

func RandStringBytesMask(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}
