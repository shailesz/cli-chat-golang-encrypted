package controllers

import (
	"fmt"

	"github.com/shailesz/cli-chat-golang/cryptoutils"
	"github.com/shailesz/cli-chat-golang/src/helpers"
	"github.com/shailesz/cli-chat-golang/src/models"
	"github.com/shailesz/cli-chat-golang/src/services"
	"github.com/shailesz/cli-chat-golang/src/socket"
)

func setupChatroom(sharedSecret []byte, config models.Config) {
	helpers.ClearScreen()
	helpers.WelcomeText()

	// event listener for message
	socket.Socket.On("message", func(chat models.ChatMessage) {
		message, err := cryptoutils.DecryptMessageAES(sharedSecret, chat)
		if err != nil {
			fmt.Println("Failed to decrypt message:", err)
			return
		}

		services.DisplayMessage(chat.Username, message.ToString(), config.Username)
	})

	// handle input for chatroom
	services.HandleChatInput(config, sharedSecret)
}
