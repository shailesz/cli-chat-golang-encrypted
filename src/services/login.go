package services

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/shailesz/cli-chat-golang/keymanager"
	"github.com/shailesz/cli-chat-golang/src/socket"
)

// getSharedSecretForUser retrieves the shared secret for communication with a given user.
func GetSharedSecretForUser(targetUsername string, localPrivateKey []byte) ([]byte, error) {
	isReceived := false
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey

	socket.Socket.Emit("getPublicKey", targetUsername)

	socket.Socket.On("getPublicKey", func(key string) {

		publicK, err := keymanager.DecodePublicKey(key)
		if err != nil {
			fmt.Println(err)
		}

		privateK, err := keymanager.DecodePrivateKey(localPrivateKey)
		if err != nil {
			fmt.Println(err)
		}

		privateKey = privateK
		publicKey = publicK
		isReceived = true
	})

	for !isReceived {
	}

	return keymanager.DeriveSharedSecret(privateKey, publicKey)
}
