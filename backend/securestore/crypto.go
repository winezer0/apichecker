package securestore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "aesgcm:"

// EncryptText 使用机器码派生出的密钥对文本进行 AES-GCM 加密。
func EncryptText(plainText, machineID string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(deriveKey(machineID))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptText 使用机器码派生出的密钥对文本进行 AES-GCM 解密。
func DecryptText(cipherText, machineID string) (string, error) {
	trimmed := strings.TrimSpace(cipherText)
	if trimmed == "" {
		return "", nil
	}
	if !strings.HasPrefix(trimmed, encryptedPrefix) {
		return trimmed, nil
	}

	rawCipherText, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(trimmed, encryptedPrefix))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(deriveKey(machineID))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(rawCipherText) < gcm.NonceSize() {
		return "", errors.New("cipher text is too short")
	}

	nonce := rawCipherText[:gcm.NonceSize()]
	encrypted := rawCipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt text failed: %w", err)
	}

	return string(plainText), nil
}

// EncryptTextWithPassword 使用密码派生出的密钥对文本进行 AES-GCM 加密。
func EncryptTextWithPassword(plainText, password string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptTextWithPassword 使用密码派生出的密钥对文本进行 AES-GCM 解密。
func DecryptTextWithPassword(cipherText, password string) (string, error) {
	trimmed := strings.TrimSpace(cipherText)
	if trimmed == "" {
		return "", nil
	}
	if !strings.HasPrefix(trimmed, encryptedPrefix) {
		return trimmed, nil
	}

	rawCipherText, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(trimmed, encryptedPrefix))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(rawCipherText) < gcm.NonceSize() {
		return "", errors.New("cipher text is too short")
	}

	nonce := rawCipherText[:gcm.NonceSize()]
	encrypted := rawCipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt text failed: %w", err)
	}

	return string(plainText), nil
}

// deriveKey 使用 SHA-256 将输入派生为 32 字节 AES 密钥。
func deriveKey(input string) []byte {
	sum := sha256.Sum256([]byte(strings.TrimSpace(input)))
	return sum[:]
}
