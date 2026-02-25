# Prevy

A minimal clipboard history manager for the terminal, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Prevy runs a lightweight daemon that watches your system clipboard in the background. When you need something you copied earlier, open the TUI and it's all there.

![Tokyo Night](https://img.shields.io/badge/theme-Tokyo%20Night-7aa2f7?style=flat-square)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-565f89?style=flat-square)

```
 ╭─ Prevy ──────────────────────────────────────────────╮
 │                                                      │
 │  1  Lorem ipsum dolor sit amet, consec...   12s ago  │
 │  2  https://github.com/charmbracelet        1m ago   │
 │ >3  func main() { fmt.Println("hello...     5m ago   │
 │  4  SELECT * FROM users WHERE id = 42      14m ago   │
 │  5  /home/user/documents/report.pdf         1h ago   │
 │                                                      │
 ╭──────────────────────────────────────────────────────╮
 │      enter copy  p paste  x clear all  q quit       │
 ╰──────────────────────────────────────────────────────╯
```

---

## Install

### Arch Linux (AUR)

```bash
yay -S prevy
```

Or with any other AUR helper (`paru`, `trizen`, etc.), or manually:

```bash
git clone https://aur.archlinux.org/prevy.git
cd prevy
makepkg -si
```

### Go install

```bash
go install github.com/atiladefreitas/prevy@latest
```

### Build from source

```bash
git clone https://github.com/atiladefreitas/prevy.git
cd prevy
go build -o prevy .
```

### Requirements

- **Linux**: `wl-copy`/`wl-paste` (Wayland) or `xclip`/`xsel` (X11)
- **macOS**: works out of the box (`pbcopy`/`pbpaste`)

---

## Usage

### 1. Start the daemon

The daemon watches your clipboard and records every new copy to disk.

```bash
prevy --daemon &
```

> Add this to your shell profile (`~/.zshrc`, `~/.bashrc`) or set up a systemd
> service so it starts automatically on login.

### 2. Open the TUI

```bash
prevy
```

That's it. Your full clipboard history is right there.

---

## Keybindings

| Key | Action |
|---|---|
| `j` / `Down` | Move down |
| `k` / `Up` | Move up |
| `Enter` | Copy selected item to clipboard |
| `p` | Copy to clipboard and paste to stdout |
| `x` | Clear all history |
| `q` / `Esc` | Quit |

---

## Flags

| Flag | Description |
|---|---|
| *(none)* | Open the clipboard history TUI |
| `--daemon` | Start the background clipboard watcher |
| `--status` | Check if the daemon is running |
| `--version` | Show version |
| `--help` | Show help |

---

## How it works

```
┌──────────┐    polls every 1s    ┌───────────────────────┐
│  System   │ ──────────────────> │    prevy --daemon     │
│ Clipboard │                     │  (background process) │
└──────────┘                      └───────────┬───────────┘
                                              │
                                        saves to disk
                                              │
                                              v
                                  ┌───────────────────────┐
                                  │  ~/.local/share/prevy  │
                                  │     history.json       │
                                  └───────────┬───────────┘
                                              │
                                         reads on launch
                                              │
                                              v
                                  ┌───────────────────────┐
                                  │       prevy (TUI)      │
                                  │    Bubble Tea + Lip    │
                                  │        Gloss           │
                                  └───────────────────────┘
```

- The **daemon** polls the system clipboard every second. When it detects new content, it deduplicates and saves to `~/.local/share/prevy/history.json`.
- The **TUI** reads that file on launch and renders the list. Selecting an item writes it back to the system clipboard.
- History is capped at **100 entries**.

---

## Project Structure

```
prevy/
  main.go              Entry point and CLI flag handling
  clipboard/
    clipboard.go       System clipboard read/write bridge
  daemon/
    daemon.go          Background clipboard watcher with PID management
  store/
    store.go           JSON persistence and history operations
  ui/
    model.go           Bubble Tea model (Init, Update, View)
    styles.go          Lip Gloss styles (Tokyo Night palette)
    keys.go            Keymap definitions
```

---

## Autostart (systemd)

A systemd user service is included. If you installed from the AUR, just enable it:

```bash
systemctl --user enable --now prevy.service
```

If you installed via `go install`, copy the service file manually:

```bash
cp prevy.service ~/.config/systemd/user/prevy.service
# Edit ExecStart to point to your binary, e.g. ExecStart=%h/go/bin/prevy --daemon
systemctl --user daemon-reload
systemctl --user enable --now prevy.service
```

---

## License

MIT
