package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func sign(signaturePayload string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(signaturePayload))
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature
}

func ValidMAC(message string, signature string, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
