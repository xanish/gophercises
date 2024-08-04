package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt will take in a key and plaintext and return a hex representation
// of the encrypted value.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBEncrypter
func Encrypt(key, plaintext string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique but not secure. Therefore, it is common to include it at the start of ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Read random bytes into the IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt will take in a key and a cipherHex (hex representation of
// the ciphertext) and decrypt it.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBDecrypter
func Decrypt(key, cipherHex string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// newCipherBlock converts our key of any size to 16 bytes using md5.
// We need to do this as aes.NewCipher needs keys of length 16, 24 or 32 bytes specifically.
func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()

	// write the key to hasher's buffer, could also use io.WriteString or hasher.Write
	fmt.Fprint(hasher, key)

	cipherKey := hasher.Sum(nil)

	return aes.NewCipher(cipherKey)
}
