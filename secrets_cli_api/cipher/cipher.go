package cipher

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
	// The IV needs to be unique but not secure. Therefore, it is common to include it at the start of ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Read random bytes into the IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return "", err
	}

	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// EncryptWriter will return a writer that will write encrypted data to
// the original writer.
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}

	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// Decrypt will take in a key and a cipherHex (hex representation of
// the ciphertext) and decrypt it.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBDecrypter
func Decrypt(key, cipherHex string) (string, error) {
	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream, err := decryptStream(key, iv)
	if err != nil {
		return "", err
	}

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// DecryptReader will return a reader that will decrypt data from the
// provided reader and give the user a way to read that data as it if was
// not encrypted.
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to read the full iv")
	}

	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{S: stream, R: r}, nil
}

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
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
