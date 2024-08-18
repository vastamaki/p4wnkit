package menu

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/buger/goterm"
	"golang.org/x/term"
)

var (
	up    = byte(65)
	down  = byte(66)
	enter = byte(13)
	keys  = map[byte]bool{
		up:   true,
		down: true,
	}
)

var menu Menu
var exitChannel chan struct{}

type Menu struct {
	Prompt    string
	CursorPos int
	MenuItems []*MenuItem
}

type MenuItem struct {
	Text string
	ID   string
}

func NewMenu(exitChan chan struct{}) *Menu {
	exitChannel = exitChan
	menu = Menu{
		Prompt:    "",
		CursorPos: 0,
		MenuItems: make([]*MenuItem, 0),
	}

	return &menu
}

func GetMenu() *Menu {
	return &menu
}

func (m *Menu) SetPrompt(prompt string) {
	menu.Prompt = prompt
}

func (m *Menu) AskChoice(items []*MenuItem) string {
	menu.MenuItems = items
	menu.CursorPos = 0

	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()

	fmt.Println(goterm.Color(goterm.Bold(m.Prompt)+":", goterm.CYAN))

	m.renderMenuItems(false)

	for {
		keyCode := getInput()
		if keyCode == enter {
			menuItem := m.MenuItems[m.CursorPos]
			return menuItem.ID
		} else if keyCode == up {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if keyCode == down {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
	}
}

func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.CursorPos {
			cursor = goterm.Color("> ", goterm.YELLOW)
			menuItemText = goterm.Color(menuItemText, goterm.YELLOW)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

func getInput() byte {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return 0
	}

	defer func() {
		if restoreErr := term.Restore(fd, oldState); restoreErr != nil {
			fmt.Printf("Error restoring terminal state: %v\n", restoreErr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		exitChannel <- struct{}{}
	}()

	readBytes := make([]byte, 3)
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		fmt.Println("Error reading from terminal:", err)
		return 0
	}

	if readBytes[0] == 0x03 {
		exitChannel <- struct{}{}
		return 0
	}

	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}
