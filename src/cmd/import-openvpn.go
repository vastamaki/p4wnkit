package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	httpclient "github.com/vastamaki/p4wnkit/src/http"
	"github.com/vastamaki/p4wnkit/src/menu"
	"github.com/vastamaki/p4wnkit/src/utils"
)

func OpenVPN(config *utils.AppConfig) {
	items := []*menu.MenuItem{
		{
			Text: "Change VPN Server",
			ID:   "change-vpn-server",
		},
		{
			Text: "Go back",
			ID:   "back",
		},
	}

	menu := menu.GetMenu()

	menu.SetPrompt("OpenVPN Settings")

	choice := menu.AskChoice(items)

	switch choice {
	case "change-vpn-server":
		changeVPNServer(config)
	}
}

func changeVPNServer(config *utils.AppConfig) {
	defer func() {
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}()

	items := []*menu.MenuItem{
		{
			Text: "EU FREE 1",
			ID:   "1",
		},
		{
			Text: "EU FREE 2",
			ID:   "201",
		},
		{
			Text: "EU FREE 3",
			ID:   "253",
		},
		{
			Text: "Go back",
			ID:   "back",
		},
	}

	menu := menu.GetMenu()

	menu.SetPrompt("Select server")

	choice := menu.AskChoice(items)

	if choice == "back" {
		return
	}

	utils.Spinner.Start("Switching VPN server...")

	url := fmt.Sprintf("https://labs.hackthebox.com/api/v4/connections/servers/switch/%s", choice)

	_, error := httpclient.Post(url, config.HtbAccessToken)

	if error != nil {
		utils.Spinner.Error("Failed to switch VPN server")
	}

	utils.Spinner.Success("VPN server switched successfully")

	utils.Spinner.Start("Downloading OpenVPN config...")

	url = fmt.Sprintf("https://labs.hackthebox.com/api/v4/access/ovpnfile/%s/0", choice)

	response, error := httpclient.Get(url, config.HtbAccessToken)

	if error != nil {
		fmt.Print(error)
		fmt.Println("Failed to download OpenVPN config")
	}

	err := os.WriteFile(filepath.Join(config.ConfigDir, "openvpn", "vpn.ovpn"), response, 0755)
	if err != nil {
		fmt.Println(err)

		utils.Spinner.Error("Failed to save OpenVPN config")

		time.Sleep(5 * time.Second)
	}

	utils.Spinner.Success("OpenVPN config saved successfully")
}
