package clipboard

import (
	"os/exec"
	"runtime"
	"strings"
)

type provider struct {
	readCmd   string
	readArgs  []string
	writeCmd  string
	writeArgs []string
}

func detect() *provider {
	switch runtime.GOOS {
	case "darwin":
		return &provider{
			readCmd: "pbpaste", readArgs: nil,
			writeCmd: "pbcopy", writeArgs: nil,
		}
	case "linux":
		// prefer wl-copy/wl-paste on Wayland
		if _, err := exec.LookPath("wl-paste"); err == nil {
			return &provider{
				readCmd: "wl-paste", readArgs: []string{"--no-newline"},
				writeCmd: "wl-copy", writeArgs: nil,
			}
		}
		if _, err := exec.LookPath("xclip"); err == nil {
			return &provider{
				readCmd: "xclip", readArgs: []string{"-selection", "clipboard", "-o"},
				writeCmd: "xclip", writeArgs: []string{"-selection", "clipboard"},
			}
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return &provider{
				readCmd: "xsel", readArgs: []string{"--clipboard", "--output"},
				writeCmd: "xsel", writeArgs: []string{"--clipboard", "--input"},
			}
		}
	}
	return nil
}

func Read() (string, error) {
	p := detect()
	if p == nil {
		return "", nil
	}
	out, err := exec.Command(p.readCmd, p.readArgs...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func Write(text string) error {
	p := detect()
	if p == nil {
		return nil
	}
	cmd := exec.Command(p.writeCmd, p.writeArgs...)
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
