package vault

import (
	"errors"
	"github.com/xanish/gophercises/secrets_cli_api/encrypt"
)

type Vault struct {
	encodingKey string
	keyValues   map[string]string
}

func Memory(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}

func (v *Vault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}

	ret, err := encrypt.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (v *Vault) Set(key, value string) error {
	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}

	v.keyValues[key] = encryptedValue

	return nil
}
