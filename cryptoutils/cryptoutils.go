package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/shailesz/cli-chat-golang/src/models"
)

func DeriveSharedSecret(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, error) {
	if privateKey.PublicKey.Curve != publicKey.Curve {
		return nil, fmt.Errorf("curve mismatch")
	}
	x, _ := privateKey.Curve.ScalarMult(publicKey.X, publicKey.Y, privateKey.D.Bytes())
	return x.Bytes(), nil
}

func hashSharedSecret(secret []byte) []byte {
	hash := sha256.Sum256(secret)
	return hash[:]
}

func DecryptMessageAES(sharedSecret []byte, message models.ChatMessage) (models.ChatMessage, error) {
	key := hashSharedSecret(sharedSecret)

	ciphertext, err := hex.DecodeString(message.Data)
	if err != nil {
		fmt.Println("Failed to decode hex:", err)
		return message, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Failed to create new cipher:", err)
		return message, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Failed to create new GCM:", err)
		return message, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return message, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("plain text fail:", err)
		return message, err
	}

	message.Data = string(plaintextBytes)

	return message, nil
}

func EncryptMessageAES(sharedSecret, plaintext []byte) (string, error) {
	key := hashSharedSecret(sharedSecret) // Ensure the key is of the correct length for AES
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}
