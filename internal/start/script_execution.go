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

type Script struct {
	colors       map[string]tcell.Color
	IsScriptMode bool
	ScriptLines  []string
	ScriptIndex  int
}

func LoadScript(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func (s *Script) ExecuteScriptOS(app *tview.Application, input *tview.InputField, output *tview.TextView) {
	for i := 0; i < len(s.ScriptLines); i++ {
		line := s.ScriptLines[s.ScriptIndex]
		s.ScriptIndex++

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			return
		}

		fmt.Fprintf(output, "$ %s\n", line)

		args := strings.Fields(line)
		if len(args) == 0 {
			return
		}

		cmd := args[0]

		switch cmd {
		case CD:
			if len(args) > 1 {
				path := args[1]
				if err := commands.CdCMD(input, path); err != nil {
					fmt.Fprintf(output, "$ %s\n", err)
				}
			} else {
				fmt.Fprintf(output, "$ No args\n")
			}
		case LS:
			if err := commands.LsCMD(output, s.colors); err != nil {
				fmt.Fprintf(output, "$ %s\n", err)
			}
		case EXIT:
			app.Stop()
			return
		case CLEAR:
			commands.ClearCMD(output)
		default:
			fmt.Fprintf(output, "> Wrong input try another command: %s\n", cmd)
		}
	}
}

func (s *Script) ExecuteScriptVFS(currentVFS *vfs.VFS, app *tview.Application, input *tview.InputField, output *tview.TextView, colors map[string]tcell.Color) {
	for i := 0; i < len(s.ScriptLines); i++ {
		line := s.ScriptLines[s.ScriptIndex]
		s.ScriptIndex++

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			return
		}

		fmt.Fprintf(output, "$ %s\n", line)

		args := strings.Fields(line)
		if len(args) == 0 {
			return
		}

		cmd := args[0]

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
	}
}
