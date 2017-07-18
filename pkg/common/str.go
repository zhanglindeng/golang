package common

import (
  "crypto/hmac"
  "crypto/sha256"
)

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
// @link https://golang.org/pkg/crypto/hmac/
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
