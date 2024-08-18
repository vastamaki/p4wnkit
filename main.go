package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/buger/goterm"
	"github.com/vastamaki/p4wnkit/src/cmd"
	"github.com/vastamaki/p4wnkit/src/menu"
	"github.com/vastamaki/p4wnkit/src/utils"
)

var config utils.AppConfig

func init() {
	utils.Spinner.Start("Initializing p4wnkit")
	config = *utils.InitializeConfig()
}

func showMenu(exitChan chan struct{}) {
	items := []*menu.MenuItem{
		{
			Text: "OpenVPN Settings",
			ID:   "openvpn",
		},
		{
			Text: "Set HTB access token",
			ID:   "set-htb-access-token",
		},
		{
			Text: "Exit",
			ID:   "exit",
		},
	}

	menu := menu.NewMenu(exitChan)

	defer func() {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Flush()
	}()

	for {
		fmt.Print("\033[?25l")
		menu.SetPrompt("Welcome to p4wnkit, what shall we do?")

		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Flush()

		choice := menu.AskChoice(items)

		switch choice {
		case "openvpn":
			cmd.OpenVPN(&config)
		case "set-htb-access-token":
			cmd.SetHtbAccessToken(&config)
		case "exit":
			exitChan <- struct{}{}
		}
	}
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	exitChan := make(chan struct{})

	utils.Spinner.Success("p4wnkit initialized")

	go showMenu(exitChan)

	select {
	case <-done:
	case <-exitChan:
	}

	os.Exit(0)
}
