package phpgo

import (
	"testing"
)

func TestAesDecrypt(t *testing.T) {
	key := []byte("e9rnM5YDQrGbUsjx23nrFgHyrML6xvWk")
	s := "QdLbEmUjgm9N7jQ1DsMWUBMQvoPSvpa8N3LbDBswCKw"
	text, err := AesDecrypt(key, s)
	if err != nil {
		t.Error("AesDecrypt Error", err)
		return
	}
	t.Log("AesDecrypt OK", text)
}

func TestAesEncrypt(t *testing.T) {
	key := []byte("e9rnM5YDQrGbUsjx23nrFgHyrML6xvWk")
	text := "123456"
	s, err := AesEncrypt(key, text)
	if err != nil {
		t.Error("AesEncrypt Error", err)
		return
	}
	t.Log("AesEncrypt OK", s)
}
