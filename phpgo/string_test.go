package phpgo

import (
	"testing"
)

func TestStringRandomStr(t *testing.T) {
	n := 6
	s := StringRandomStr(n, []byte(`abcdefghijklmnopqrstuvwxyz`)...)

	if len(s) != n {
		t.Error("StringRandomStr Error")
	}

	t.Logf("%s", s)
}

func TestStringRandomByte(t *testing.T) {
	n := 16
	b := StringRandomByte(n, []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()_-=+|{}[];.,'"<>?~\/`)...)

	if len(b) != n {
		t.Error("StringRandomByte Error")
	}

	t.Logf("%s", b)
}

func TestStringPasswordHash(t *testing.T) {
	password := "123456"
	// ignore error for the sake of simplicity
	hash, _ := StringPasswordHash(password)

	t.Log("Password:", password)
	t.Log("Hash:    ", hash)

	match := StringPasswordVerify(password, hash)
	if !match {
		t.Error("StringPasswordHash Error: Not match")
	}
}

func TestStringSha256(t *testing.T) {
	s := "123456"
	result := StringSha256(s)
	// 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
	if len(result) != 64 {
		t.Error("StringSha256 Error")
	}
	t.Logf("StringSha256: %s => %s", s, result)
}

func TestStringSha1(t *testing.T) {
	s := "123456"
	result := StringSha1(s)
	// 7c4a8d09ca3762af61e59520943dc26494f8941b
	if len(result) != 40 {
		t.Error("StringSha1 Error")
	}
	t.Logf("StringSha1: %s => %s", s, result)
}

func TestStringMd5(t *testing.T) {
	s := "123456"
	result := StringMd5(s)
	// expected := "e10adc3949ba59abbe56e057f20f883e"
	if len(result) != 32 {
		t.Error("StringMd5 Error")
	}
	t.Logf("StringMd5: %s => %s", s, result)
}
