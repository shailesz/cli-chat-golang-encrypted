package controllers

import (
	"fmt"

	"github.com/shailesz/cli-chat-golang/keymanager"
	"github.com/shailesz/cli-chat-golang/src/helpers"
	"github.com/shailesz/cli-chat-golang/src/models"
	"github.com/shailesz/cli-chat-golang/src/services"
	"github.com/shailesz/cli-chat-golang/src/socket"
)

func HandleLogin(username, receiver string) error {
	c := helpers.ReadConfig(username)

	key, err := keymanager.DecryptPrivateKey(c.EncryptedPrivateKey, c.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	sharedSecret, err := services.GetSharedSecretForUser(receiver, key)
	if err != nil {
		return fmt.Errorf("error getting shared secret: %w", err)
	}

	var isWaiting, isUpdate bool
	var u, p, pk string

	// handle configs from config file
	if c.Username == "" || c.Password == "" {
		_, u, p = helpers.GetCredentials(false)

		isUpdate = true

	} else {
		fmt.Println("Processing...")
		u, p, pk = c.Username, c.Password, c.PublicKey
		isUpdate = false
	}

	// listener for auth messages.
	socket.Socket.On("auth", func(message models.AuthMessage) {
		if message.Status == 404 {
			fmt.Println("You could not be authenticated. please try again.")
		} else {
			fmt.Println("Authenticated.")

			if isUpdate {
				c.Update(u, p)
			}
		}

		isWaiting = false
	})

	isWaiting = true
	socket.Socket.Emit("auth", models.User{Username: u, Password: p, PublicKey: pk}) // send auth message to server

	// wait for auth message.
	for isWaiting {
	}
	setupChatroom(sharedSecret, c)

	return nil
}
