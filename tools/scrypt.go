package tools

import (
	"bytes"
	r "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"strings"
	"time"
)

// RSA加密
func RSAEncrypt(data, key []byte) []byte {
	block, _ := pem.Decode(key)
	if block == nil {
		zap.S().Error("block key fail")
		return nil
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		zap.S().Error("decrypt key fail", err)
		return nil
	}
	pub := pubInterface.(*rsa.PublicKey)
	text, err := rsa.EncryptPKCS1v15(r.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return text
}

// RSA 解密
func RSADecrypt(data, key []byte) []byte {
	block, _ := pem.Decode(key)
	if block == nil {
		zap.S().Error("block key fail")
		return nil
	}
	prev, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		zap.S().Error(err)
		return nil
	}
	data, err = rsa.DecryptPKCS1v15(r.Reader, prev, data)
	if err != nil {
		zap.S().Error(err)
	}
	return data
}

// 生成RSA
func GenRSAEncrypt() ([]byte, []byte, error) {
	private, err := rsa.GenerateKey(r.Reader, 1024)
	if err != nil {
		zap.S().Error(err)
		return nil, nil, errors.New("Gen RSA Fail ")
	}
	derStream := x509.MarshalPKCS1PrivateKey(private)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateKey := pem.EncodeToMemory(block)
	public := &private.PublicKey
	derPix, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPix,
	}
	publicKey := pem.EncodeToMemory(block)
	return privateKey, publicKey, nil
}

// 校验明文+盐 是否与 code 相等
func Compare(text, salt, code string) bool {
	if len(text) == 0 || len(salt) == 0 || len(code) == 0 {
		return false
	}
	dk, err := scrypt.Key([]byte(text), []byte(salt), 1024, 8, 1, 32)
	if err != nil {
		return false
	}
	return strings.Compare(code, hex.EncodeToString(dk)) == 0
}

// 生成密文code
func GenCode(text string) (string, string, error) {
	salt := getRandomSalt()
	dk, err := scrypt.Key([]byte(text), []byte(salt), 1024, 8, 1, 32)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(dk), salt, nil
}

// 随机盐种子
var seed = [86]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"~", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "+", ",", ".", "/", ";", ":", "'", "\\", "|", "`", "<", ">", "?",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

// 生成随机盐
func getRandomSalt() string {
	var salt bytes.Buffer
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 10; i++ {
		salt.WriteString(seed[rand.Intn(len(seed))])
	}
	return salt.String()
}
