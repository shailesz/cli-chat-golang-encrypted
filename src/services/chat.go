package services

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/shailesz/cli-chat-golang/cryptoutils"
	"github.com/shailesz/cli-chat-golang/src/constants"
	"github.com/shailesz/cli-chat-golang/src/helpers"
	"github.com/shailesz/cli-chat-golang/src/models"
	"github.com/shailesz/cli-chat-golang/src/socket"
)

func DisplayMessage(sender, message, currentUser string) {
	if sender != currentUser {
		helpers.ClearLine()
		fmt.Println(constants.PURPLE_TERMINAL_COLOR + message + constants.RESET_TERMINAL_COLOR)
		helpers.Prompt()
	}
}

// SendChat emits chat event to server.
func SendChat(u, m string, sharedSecret []byte) {
	encryptedMessage, err := cryptoutils.EncryptMessageAES(sharedSecret, []byte(m))
	if err != nil {
		fmt.Println("Failed to encrypt message:", err)
		return
	}

	socket.Socket.Emit("chat", models.ChatMessage{Username: u, Data: encryptedMessage, Timestamp: time.Now().UnixNano()})
}

// HandleChatInput sends scanned input from to server.
func HandleChatInput(config models.Config, sharedSecret []byte) {
	reader := bufio.NewReader(os.Stdin)

	// prompt
	for {
		helpers.Prompt()
		data, _, _ := reader.ReadLine()
		message := string(data)
		SendChat(config.Username, message, sharedSecret)
		if message == "$quit" {
			break
		}
	}
}
