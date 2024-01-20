package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

// Aes crypto
type Aes struct {
	secret string
}

// New aes
//
// secret length must be 16, 24, or 32 bytes
func NewAes(secret string) Aes {
	return Aes{secret}
}

// Encryt in CBC Mode
//
// plainText is string
//
// iv length must be same with with block size = 16 bytes = 128 bits
func (a *Aes) EncryptCBC(
	plainText string,
	iv string,
) (*string, error) {
	block, err := aes.NewCipher([]byte(a.secret))
	ivBytes := []byte(iv)

	bPlaintext := PKCS5Padding([]byte(plainText), aes.BlockSize, len(plainText))

	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(bPlaintext))
	blockMode := cipher.NewCBCEncrypter(block, ivBytes)
	blockMode.CryptBlocks(ciphertext, bPlaintext)
	result := hex.EncodeToString(ciphertext)

	return &result, nil
}

// Decrypt cipherText to plain text with CBC Mode
//
// cipherText is hex string
//
// iv length must be same with with block size = 16 bytes = 128 bits
//
// return plainText
func (a *Aes) DecryptCBC(
	cipherText string,
	iv string,
) (*string, error) {
	cipherTextBytes, err := hex.DecodeString(cipherText)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(a.secret))

	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))

	mode.CryptBlocks([]byte(cipherTextBytes), []byte(cipherTextBytes))
	result := string(PKCS5UnPadding(cipherTextBytes))

	return &result, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
