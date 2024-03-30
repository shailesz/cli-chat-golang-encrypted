package keymanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

func GenerateECDSAKeys() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, "", err
	}
	publicKeyBytes := elliptic.Marshal(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
	return privateKey, hex.EncodeToString(publicKeyBytes), nil
}

func EncryptPrivateKey(privateKey *ecdsa.PrivateKey, passphrase string) (string, error) {
	privateKeyBytes := privateKey.D.Bytes()
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New)
	block, err := aes.NewCipher(key)
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
	ciphertext := gcm.Seal(nonce, nonce, privateKeyBytes, nil)
	ciphertextWithSalt := append(salt, ciphertext...)
	return hex.EncodeToString(ciphertextWithSalt), nil
}

func DecodePublicKey(pubKeyStr string) (*ecdsa.PublicKey, error) {
	pubBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

func DecodePrivateKey(privKeyBytes []byte) (*ecdsa.PrivateKey, error) {

	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = new(big.Int).SetBytes(privKeyBytes)
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKey.D.Bytes())

	return privKey, nil
}

func DeriveSharedSecret(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, error) {
	// Perform scalar multiplication directly
	x, _ := privateKey.Curve.ScalarMult(publicKey.X, publicKey.Y, privateKey.D.Bytes())
	return x.Bytes(), nil
}

func DecryptPrivateKey(encryptedPrivateKeyHex, passphrase string) ([]byte, error) {
	encryptedData, err := hex.DecodeString(encryptedPrivateKeyHex)
	if err != nil {
		return nil, err
	}

	salt := encryptedData[:16] // Assuming the first 16 bytes are the salt
	encryptedPrivateKey := encryptedData[16:]

	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedPrivateKey) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedPrivateKey[:nonceSize], encryptedPrivateKey[nonceSize:]
	decryptedPrivateKeyBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedPrivateKeyBytes, nil
}
