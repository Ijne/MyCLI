package start

import (
	vfs "MyCLI/internal/VFS"
	"MyCLI/internal/commands"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CD     = "cd"
	LS     = "ls"
	EXIT   = "exit"
	CLEAR  = "clear"
	WHOAMI = "whoami"
	WC     = "wc"
	TOUCH  = "touch"
)

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

func handler_os(app *tview.Application, input *tview.InputField, output *tview.TextView, colors map[string]tcell.Color) {
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
					if err := commands.CdCMD(input, path); err != nil {
						fmt.Fprintf(output, "$ %s\n", err)
					}
					//} else {
					//	fmt.Fprintf(output, "$ If you wanna send an argument, use '{arg}'\n")
					//}
				} else {
					fmt.Fprintf(output, "$ No args\n")
				}
			case LS:
				if err := commands.LsCMD(output, colors); err != nil {
					fmt.Fprintf(output, "$ %s\n", err)
				}
			case WHOAMI:
				user := os.Getenv("USER")
				if user == "" {
					user = os.Getenv("USERNAME")
				}
				if user == "" {
					user = "unknown"
				}
				fmt.Fprintf(output, "> %s\n", user)
			case WC:
				if len(args) > 1 {
					filename := args[1]
					content, err := os.ReadFile(filename)
					if err != nil {
						fmt.Fprintf(output, "wc: %s: %s\n", filename, err)
					} else {
						lines := strings.Count(string(content), "\n")
						words := len(strings.Fields(string(content)))
						bytes := len(content)
						fmt.Fprintf(output, "> %d %d %d %s\n", lines, words, bytes, filename)
					}
				} else {
					fmt.Fprintf(output, "wc: missing file operand\n")
				}
			case TOUCH:
				if len(args) > 1 {
					filename := args[1]
					file, err := os.Create(filename)
					if err != nil {
						fmt.Fprintf(output, "touch: %s: %s\n", filename, err)
					} else {
						file.Close()
						fmt.Fprintf(output, "> Created: %s\n", filename)
					}
				} else {
					fmt.Fprintf(output, "touch: missing file operand\n")
				}
			case CLEAR:
				commands.ClearCMD(output)
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
}

func handler_vfs(currentVFS *vfs.VFS, app *tview.Application, input *tview.InputField, output *tview.TextView, colors map[string]tcell.Color) {
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := input.GetText()
		args := strings.Fields(text)
		if len(args) == 0 {
			input.SetText("")
			return
		}

		cmd := args[0]

		fmt.Fprintf(output, "$ %s\n", text)

		switch cmd {
		case CD:
			if len(args) > 1 {
				path := args[1]
				if err := currentVFS.CD(path); err != nil {
					fmt.Fprintf(output, "Error: %s\n", err)
				} else {
					commands.UpdateInputLabel(input)
				}
			} else {
				currentVFS.Current = currentVFS.Root
				commands.UpdateInputLabel(input)
			}

		case LS:
			var path string
			if len(args) > 1 {
				path = args[1]
			}

			files, err := currentVFS.LS(path)
			if err != nil {
				fmt.Fprintf(output, "Error: %s\n", err)
			} else {
				fmt.Fprintf(output, "> ")
				for i, file := range files {
					if i > 0 {
						fmt.Fprintf(output, "  ")
					}
					node, _ := currentVFS.FindNode(file)
					if node != nil && node.IsDir {
						fmt.Fprintf(output, "[#%06X]%s[#%06X]", colors["dir"].Hex(), file, colors["output_text"].Hex())
					} else {
						fmt.Fprintf(output, "[#%06X]%s[#%06X]", colors["file"].Hex(), file, colors["output_text"].Hex())
					}
				}
				fmt.Fprintf(output, "\n")
			}

		case "whoami":
			fmt.Fprintf(output, "> %s\n", currentVFS.Username)

		case "wc":
			if len(args) > 1 {
				content, err := currentVFS.GetContent(args[1])
				if err != nil {
					fmt.Fprintf(output, "Error: %s\n", err)
				} else {
					lines := strings.Count(content, "\n") + 1
					words := len(strings.Fields(content))
					bytes := len(content)
					fmt.Fprintf(output, "> %d %d %d %s\n", lines, words, bytes, args[1])
				}
			} else {
				fmt.Fprintf(output, "Error: wc requires filename argument\n")
			}

		case "touch":
			if len(args) > 1 {
				if err := currentVFS.Touch(args[1]); err != nil {
					fmt.Fprintf(output, "Error: %s\n", err)
				} else {
					fmt.Fprintf(output, "> Created file: %s\n", args[1])
				}
			} else {
				fmt.Fprintf(output, "Error: touch requires filename argument\n")
			}

		case CLEAR:
			output.SetText("")

		case EXIT:
			app.Stop()
			return

		default:
			fmt.Fprintf(output, "> Unknown command: %s\n", cmd)
			fmt.Fprintf(output, "> Available commands: cd, ls, whoami, wc, touch, clear, exit\n")
		}

		input.SetText("")
		output.ScrollToEnd()
	})
}

func StartAPP() {
	config := parseCommandLineFlags()

	fmt.Fprintf(os.Stderr, "[DEBUG] Config parameters:\n")
	fmt.Fprintf(os.Stderr, "[DEBUG]   VFS Path: %s\n", config.VFSPath)
	fmt.Fprintf(os.Stderr, "[DEBUG]   Script Path: %s\n", config.ScriptPath)

	colors := chooseTheme()

	app := tview.NewApplication()

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
	commands.UpdateInputLabel(input)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("MyCLI").SetBorder(true)
	flex.AddItem(output, 0, 1, false).AddItem(input, 1, 0, true)

	var isVFS bool = false
	vfs := vfs.NewVFS()
	if config.VFSPath != "" {
		if err := vfs.LoadFromCSV(config.VFSPath); err != nil {
			fmt.Fprintf(output, "Error loading VFS: %v\n Started in interacive mode\n", err)
			handler_os(app, input, output, colors)
		} else {
			handler_vfs(vfs, app, input, output, colors)
			isVFS = true
		}
	} else {
		handler_os(app, input, output, colors)
	}
	if config.ScriptPath != "" {
		scriptLines, err := LoadScript(config.ScriptPath)
		if err != nil {
			fmt.Fprintf(output, "Error loading script: %v\n Started in interacive mode\n", err)
		} else {
			input.SetText(" ▁▂▃▅▇    SYSTEM: To exit application, please press `Ctrl + C`    ▇▅▃▂▁")
			input.SetDisabled(true)
		}

		s := Script{
			colors:       colors,
			IsScriptMode: true,
			ScriptLines:  scriptLines,
		}

		go func() {
			app.QueueUpdateDraw(func() {
				if isVFS {
					s.ExecuteScriptVFS(vfs, app, input, output, colors)
				} else {
					s.ExecuteScriptOS(app, input, output)
				}
			})
		}()
	}

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
