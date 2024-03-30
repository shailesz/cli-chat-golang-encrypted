# cli-chat-golang
Client side cli application to chat with people who are online and connected to the server.

Built using golang <3

## Usage

To run the app: `go run main.go` 

### Available commands: 

`signup`: saves a config file in local client storage and save required credentials in db

`login <you_username> <receiver_username>`: log in using the config saved during signup and send to the receiver username if the user exists in db. 

NOTE: receiver must be signed up before you can send them messages, proper error handling has not been done in this example.

## Screenshots
![test](https://github.com/shailesz/cli-chat-golang-encrypted/assets/40053781/6f8e2250-337b-4bfc-84f7-18d707c263f0)
> client cli

![db-sc](https://github.com/shailesz/cli-chat-golang-encrypted/assets/40053781/556f337f-5cee-43a6-80c3-60a236ac930c)
> messages encrypted in the db that require key pairs, salt and passphrase to decrypt

