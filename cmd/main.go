package main

import (
	"fmt"
	"os"
	
	"github.com/charmbracelet/bubbletea"
	"github.com/kidyme/git-msg-picker/internal/picker"
)

func main() {
	p := tea.NewProgram(picker.GetPickerModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
