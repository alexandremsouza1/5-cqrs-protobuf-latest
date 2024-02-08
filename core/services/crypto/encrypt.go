package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"main.go/services/logger"
)

// Encrypt will encrypt a raw string to
// an encrypted value
// an encrypted value has an IV (nonce) + actual encrypted value
// when we decrypt, we only decrypt the latter part
func Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	secretKey, err := getSecret(ctx)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(iv, iv, data, nil)

	return ciphertext, nil
}

func EncryptHEX(ctx context.Context, data []byte) (*string, error) {
	secretKey, err := getSecret(ctx)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(iv, iv, data, nil)
	ciphertextHex := hex.EncodeToString(ciphertext)

	return &ciphertextHex, nil
}

func Decrypt(ctx context.Context, data string) ([]byte, error) {
	secretKey, err := getSecret(ctx)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, cryptoText := data[:nonceSize], data[nonceSize:]

	ciphertext, err := gcm.Open(nil, []byte(nonce), []byte(cryptoText), nil)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func DecryptHEX(ctx context.Context, data string) ([]byte, error) {
	dataFromHex, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	secretKey, err := getSecret(ctx)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, cryptoText := dataFromHex[:nonceSize], dataFromHex[nonceSize:]

	ciphertext, err := gcm.Open(nil, []byte(nonce), []byte(cryptoText), nil)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func getSecret(ctx context.Context) ([]byte, error) {
	var err error
	secret := os.Getenv("CORE_AES_SECRET")
	secretbite, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		// probably malform secret, panic out
		logger.Error(ctx, err)
	}
	return secretbite, err
}

func NewAESRandomKey() error {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// really, what are you gonna do if randomness failed?
		return err
	}
	base64Key := base64.StdEncoding.EncodeToString([]byte(key))
	fmt.Printf("AES Secret:       %s\n", base64Key)
	return nil
}

func NewRandomKey(ctx context.Context) []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// really, what are you gonna do if randomness failed?
		logger.Error(ctx, err)
	}
	return key
}
