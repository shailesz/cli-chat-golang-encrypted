package cmd

import (
	"fmt"

	"github.com/shailesz/cli-chat-golang/src/controllers" // Assuming services package exists
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [username] [receiver]",
	Short: "Log in to the chat application",
	Long: `Log in to the chat application using your username and start a chat with the receiver.
Example usage:
	chatapp login john jane`,
	Args: cobra.ExactArgs(2), // Ensures exactly two arguments are passed
	Run:  loginRun,
}

func loginRun(cmd *cobra.Command, args []string) {
	login, receiver := args[0], args[1]

	if err := controllers.HandleLogin(login, receiver); err != nil {
		fmt.Println("Login error:", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
