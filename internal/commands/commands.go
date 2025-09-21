package commands

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CD    = "cd"
	LS    = "ls"
	EXIT  = "exit"
	CLEAR = "clear"
)

func LsCMD(output *tview.TextView, colors map[string]tcell.Color) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	arr, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "> ")
	for ind, d := range arr {
		if (ind+1)%6 == 0 {
			fmt.Fprintf(output, "\n")
		}
		if d.IsDir() {
			fmt.Fprintf(output, "'[#%06X]%s[#%06X]' ",
				colors["dir"].Hex(),
				d.Name(),
				colors["output_text"].Hex())
		} else {
			fmt.Fprintf(output, "'[#%06X]%s[#%06X]' ",
				colors["file"].Hex(),
				d.Name(),
				colors["output_text"].Hex())
		}
	}
	fmt.Fprintf(output, "\n")

	return nil
}

func CdCMD(input *tview.InputField, path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}

	UpdateInputLabel(input)
	return nil
}

func ClearCMD(output *tview.TextView) {
	output.SetText("")
}

func getCurDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, err
}

func UpdateInputLabel(input *tview.InputField) {
	label, err := getCurDir()
	if err != nil {
		label = "$> "
	}
	input.SetLabel(label + "> ")
}
