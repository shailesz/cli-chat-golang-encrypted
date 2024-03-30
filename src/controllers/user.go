package controllers

import (
	"fmt"

	"github.com/shailesz/cli-chat-golang/keymanager"
	"github.com/shailesz/cli-chat-golang/src/helpers"
	"github.com/shailesz/cli-chat-golang/src/models"
	"github.com/shailesz/cli-chat-golang/src/socket"
)

// CreateUser creates a user.
func CreateUser(e, u, p string) {
	var waitResponse bool = true

	privateKey, publicKeyHex, err := keymanager.GenerateECDSAKeys()
	if err != nil {
		fmt.Println("Failed to generate keys:", err)
		return
	}

	passphrase := helpers.Sha256(p) // Example: using a hash of the password as the encryption passphrase
	encryptedPrivateKey, err := keymanager.EncryptPrivateKey(privateKey, passphrase)
	if err != nil {
		fmt.Println("Failed to encrypt private key:", err)
		return
	}

	user := models.User{Email: e, Username: u, Password: passphrase, PublicKey: publicKeyHex}

	// Save user configuration including the encrypted private key
	config := models.Config{
		Email:               e,
		Username:            u,
		Password:            passphrase,
		PublicKey:           publicKeyHex,
		EncryptedPrivateKey: encryptedPrivateKey,
		PrivateKey:          privateKey.D.Bytes(),
	}
	models.WriteConfig(config, u)

	// Emit signup event with user data (excluding the private key)
	socket.Socket.Emit("signup", user)

	socket.Socket.On("signup", func(res models.AuthMessage) {
		if res.Status == 200 {
			fmt.Println("Successfully signed up, please continue to login.", passphrase)
		} else {
			// Handle different statuses
		}
		waitResponse = false
	})

	for waitResponse {
		// Block until response is received
	}
}
