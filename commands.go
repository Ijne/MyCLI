package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CD    = "cd"
	LS    = "ls"
	EXIT  = "exit"
	CLEAR = "clear"
)

func lsCMD(output *tview.TextView, colors map[string]tcell.Color) error {
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

func cdCMD(input *tview.InputField, path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}

	updateInputLabel(input)
	return nil
}

func clearCMD(output *tview.TextView) {
	output.SetText("")
}

// Вспомогательные
func getCurDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, err
}

func updateInputLabel(input *tview.InputField) {
	label, err := getCurDir()
	if err != nil {
		label = "$> "
	}
	input.SetLabel(label + "> ")
}

func chooseTheme() map[string]tcell.Color {
	var theme string
	for {
		fmt.Println("Выберите тему: light, dark, contrast:")
		fmt.Scan(&theme)
		switch theme {
		case "light":
			return map[string]tcell.Color{
				"background":  tcell.NewRGBColor(245, 245, 245),
				"input_bg":    tcell.NewRGBColor(255, 255, 255),
				"input_text":  tcell.NewRGBColor(40, 40, 40),
				"output_bg":   tcell.NewRGBColor(255, 255, 255),
				"output_text": tcell.NewRGBColor(60, 60, 60),
				"dir":         tcell.NewRGBColor(30, 144, 255),
				"file":        tcell.NewRGBColor(0, 0, 0),
			}
		case "dark":
			return map[string]tcell.Color{
				"background":  tcell.NewRGBColor(30, 30, 40),
				"input_bg":    tcell.NewRGBColor(50, 50, 60),
				"input_text":  tcell.NewRGBColor(255, 255, 255),
				"output_bg":   tcell.NewRGBColor(40, 40, 50),
				"output_text": tcell.NewRGBColor(220, 220, 220),
				"dir":         tcell.NewRGBColor(135, 206, 235),
				"file":        tcell.NewRGBColor(255, 255, 255),
			}
		case "contrast":
			return map[string]tcell.Color{
				"background":  tcell.NewRGBColor(0, 0, 0),
				"input_bg":    tcell.NewRGBColor(0, 0, 0),
				"input_text":  tcell.NewRGBColor(255, 255, 255),
				"output_bg":   tcell.NewRGBColor(0, 0, 0),
				"output_text": tcell.NewRGBColor(255, 255, 255),
				"dir":         tcell.NewRGBColor(0, 255, 255),
				"file":        tcell.NewRGBColor(255, 255, 255),
			}

		default:
			fmt.Println("Выберите среди предложенных.")
		}
	}
}

func StartAPP() {
	colors := chooseTheme()

	app := tview.NewApplication()

	// Создаем компоненты
	output := tview.NewTextView().SetDynamicColors(true).SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	output.SetBackgroundColor(colors["output_bg"])
	output.SetTextColor(colors["output_text"])

	input := tview.NewInputField().SetFieldWidth(0)
	input.SetLabelColor(colors["input_text"])
	input.SetBackgroundColor(colors["input_bg"])
	input.SetFieldTextColor(colors["input_text"])
	input.SetFieldBackgroundColor(colors["input_bg"])
	updateInputLabel(input)

	// Настраиваем layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("MyCLI").SetBorder(true)
	flex.AddItem(output, 0, 1, false).AddItem(input, 1, 0, true)

	// Обработчик
	input.SetDoneFunc(func(key tcell.Key) {
		args := strings.Fields(input.GetText())
		cmd := args[0]
		if key == tcell.KeyEnter {
			switch cmd {
			case CD:
				if len(args) > 1 {
					path := args[1]
					//if path[0] == '\'' && path[len(path)-1] == '\'' {
					//	path := path[1 : len(path)-1]
					if err := cdCMD(input, path); err != nil {
						fmt.Fprintf(output, "$ %s\n", err)
					}
					//} else {
					//	fmt.Fprintf(output, "$ If you wanna send an argument, use '{arg}'\n")
					//}
				} else {
					fmt.Fprintf(output, "$ No args\n")
				}
			case LS:
				if err := lsCMD(output, colors); err != nil {
					fmt.Fprintf(output, "$ %s\n", err)
				}
			case CLEAR:
				clearCMD(output)
			case EXIT:
				fmt.Fprintf(output, "$ %s\n", cmd)
				app.Stop()
			default:
				fmt.Fprintf(output, "> Wrong input try another command: %s\n", cmd)
			}
			input.SetText("")
			output.ScrollToEnd()
		}
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
