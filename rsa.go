package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"crypto/sha256"
	"crypto"
	"crypto/rand"
	"log"
	"encoding/pem"
	"crypto/x509"
	"io/ioutil"
	"errors"
)

func main5() {
	filename := "public.pem"
	key, err := getPublicKey(filename)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(key)
}

func main4() {
	filename := "private.pem"
	key, err := getPrivateKey(filename)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(key)
}

func main3() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := &privateKey.PublicKey
	saveKey(privateKey, publicKey)
}

func savePublicKey(key *rsa.PublicKey, filename string) error {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: b,
	})
	return ioutil.WriteFile(filename, pubBytes, 0644)
}

func savePrivateKey(key *rsa.PrivateKey, filename string) error {
	b := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return ioutil.WriteFile(filename, b, 0644)
}

func saveKey(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	ioutil.WriteFile("private.pem", privBytes, 0644)

	PubASN1, _ := x509.MarshalPKIXPublicKey(publicKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})
	ioutil.WriteFile("public.pem", pubBytes, 0644)
}

func getPublicKey(filename string) (*rsa.PublicKey, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(b)
	if p == nil {
		return nil, errors.New("invalid rsa public key")
	}

	if got, want := p.Type, "RSA PUBLIC KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("unknown key type %q, want %q", got, want))
	}
	key, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func getPrivateKey(filename string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(b)
	if p == nil {
		return nil, errors.New("invalid rsa private key")
	}

	if got, want := p.Type, "RSA PRIVATE KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("unknown key type %q, want %q", got, want))
	}
	key, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func main() {

	// Generate RSA Keys
	miryanPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	miryanPublicKey := &miryanPrivateKey.PublicKey

	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	raulPublicKey := &raulPrivateKey.PublicKey

	fmt.Println("Private Key : ", miryanPrivateKey)
	fmt.Println("Public key ", miryanPublicKey)
	fmt.Println("Private Key : ", raulPrivateKey)
	fmt.Println("Public key ", raulPublicKey)

	//Encrypt Miryan Message
	message := []byte("the code must be like a piece of music")
	label := []byte("")
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, message, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, miryanPrivateKey, newhash, hashed, &opts)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("PSS Signature : %x\n", signature)

	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, raulPrivateKey, ciphertext, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, plainText)

	//Verify Signature
	err = rsa.VerifyPSS(miryanPublicKey, newhash, hashed, signature, &opts)

	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful")
	}

}
