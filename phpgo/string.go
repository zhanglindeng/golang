package phpgo

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	r "math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// StringRandomStr generate randmon string by specify chars
func StringRandomStr(n int, alphabets ...byte) string {
	return string(StringRandomByte(n, alphabets...))
}

// StringRandomByte generate random []byte by specify chars
func StringRandomByte(n int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	}
	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = alphabets[r.Intn(len(alphabets))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes
}

// StringPasswordVerify password verfiy
func StringPasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// StringPasswordHash password hash
func StringPasswordHash(s string) (string, error) {
	// cost 设置12就够了，太大需要的计算时间太长，PHP中默认是10
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 12)
	return string(bytes), err
}

// StringSha256 sha256
func StringSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

// StringSha1 sha1
func StringSha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// StringMd5 md5
func StringMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
