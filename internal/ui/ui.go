package ui

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"teletype/internal/protocol"
)

// ANSI codes
const (
	ClearScreen = "\033[2J"
	MoveToTop   = "\033[H"
	ColorGreen  = "\033[32m"
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[34m"
	CarriageRet = "\r"
	ClearLine   = "\033[2K"
)

type UI struct {
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewUI() *UI {
	return &UI{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (ui *UI) Init() {
	fmt.Print(ClearScreen)
	fmt.Print(MoveToTop)
	fmt.Println(ColorGreen + "Welcome to TeleType... Connecting to server..." + ColorReset)
}

func (ui *UI) PrintMessage(msg protocol.Message) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	// Clear current input line
	fmt.Print(CarriageRet + ClearLine)

	timestamp := msg.Timestamp.Format("15:04:05")

	switch msg.Type {
	case protocol.MsgTypeChat:
		fmt.Printf("[%s] %s%s%s: %s\n", timestamp, ColorBlue, msg.Sender, ColorReset, msg.Content)
	case protocol.MsgTypeSystem:
		fmt.Printf("[%s] << %s >>\n", timestamp, msg.Content)
	case protocol.MsgTypeJoin:
		fmt.Printf("[%s] ++ %s joined %s ++\n", timestamp, msg.Sender, msg.Room)
	}

	// Redraw input prompt
	fmt.Print(ColorGreen + "> " + ColorReset)
}

func (ui *UI) Prompt() string {
	ui.mu.Lock()
	fmt.Print(ColorGreen + "> " + ColorReset)
	ui.mu.Unlock()

	if ui.scanner.Scan() {
		return ui.scanner.Text()
	}
	return ""
}
