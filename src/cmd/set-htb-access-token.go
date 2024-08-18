package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/buger/goterm"
	"github.com/vastamaki/p4wnkit/src/utils"
)

func SetHtbAccessToken(config *utils.AppConfig) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your HTB access token: ")
	htbAccessToken, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	config.HtbAccessToken = htbAccessToken

	if err := utils.SaveConfig(config); err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
}
