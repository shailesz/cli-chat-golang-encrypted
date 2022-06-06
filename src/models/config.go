package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/shailesz/cli-chat-golang/src/controllers"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// Config type for config file.
type Config User

// Init initializes config file.
func (c *Config) Init() Config {
	return Config{Username: "", Password: ""}
}

// Update updates config file with given parameters.
func (c *Config) Update(u, p string) Config {

	c.Username, c.Password = u, p

	// write config file.
	file, _ := json.MarshalIndent(c, "", " ")

	_ = ioutil.WriteFile("config.json", file, 0644)

	return Config{Username: u, Password: p}
}

// Login logs in user from config file.
func (c *Config) Login(socket *socketio_client.Client) (string, string) {
	var isWaiting, isUpdate bool
	var u, p string
	var err error

	// handle configs from config file
	if c.Username == "" || c.Password == "" {
		fmt.Println("Please enter credentials to continue.")
		u, p, err = controllers.Credentials()

		isUpdate = true

		if err != nil {
			log.Panicln(err)
		}
	} else {
		fmt.Println("Processing...")
		u, p = c.Username, c.Password
		isUpdate = false
	}

	// listener for auth messages.
	socket.On("auth", func(message AuthMessage) {
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
	socket.Emit("auth", User{Username: u, Password: p})

	// wait for auth message.
	for {
		if !isWaiting {
			break
		}
	}

	return u, p
}
