package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var (
	localKey = []byte("SCkGDV0EQmce7tT79hVekRiAVuBet9Ll") // AES-256 is 32 bytes
	localNonce = []byte("ga2A0X30509s") // defautlt length of 12 bytes
)

// encryptCredentials is used to encrypt locally stored credentials (eg. Deezer username and password).
func encryptCredentials(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(localNonce, localNonce, plaintext, nil), nil
}

// decryptCredentials is used to decrypt locally stored credentials (eg. Deezer username and password).
func decryptCredentials(ciphertext []byte, key[]byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
