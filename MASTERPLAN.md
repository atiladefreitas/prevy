# Prevy - Clipboard History Manager

## Overview

Prevy is a minimal, fast TUI clipboard history manager built with Go and
[Charm Bubble Tea](https://github.com/charmbracelet/bubbletea). It runs in the
terminal, displays the full clipboard history, and lets you act on it with
single-key commands.

Invoke it by typing:

```
prevy
```

---

## Design Philosophy

- **Minimal** -- no chrome, no noise; content takes center stage.
- **Fast** -- sub-100ms startup; clipboard reads are async.
- **Tokyo Night palette** -- every color pulled straight from the Tokyo Night
  theme so it feels native to developers who already use it.

---

## Tokyo Night Color Palette

| Role            | Hex       | Usage                              |
|-----------------|-----------|------------------------------------|
| Background      | `#1a1b26` | Main background                    |
| Surface         | `#24283b` | Selected item / highlighted row    |
| Overlay         | `#414868` | Borders, subtle separators         |
| Foreground      | `#c0caf5` | Primary text                       |
| Muted           | `#565f89` | Secondary text, help bar, indices  |
| Blue            | `#7aa2f7` | Accent, cursor indicator           |
| Cyan            | `#7dcfff` | Timestamps, metadata               |
| Green           | `#9ece6a` | Success messages                   |
| Red             | `#f7768e` | Destructive action highlights      |
| Yellow          | `#e0af68` | Warnings                           |

---

## TUI Layout

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
 │  enter copy to clipboard  x clear all  q quit       │
 ╰──────────────────────────────────────────────────────╯
```

### Breakdown

- **Title bar** -- app name rendered in the top border.
- **List area** -- scrollable list of clipboard entries. Each row shows:
  - Index number (muted).
  - Truncated content preview (foreground).
  - Relative timestamp (cyan).
- **Cursor** -- `>` marker + Surface background highlight on the active row.
- **Help bar** -- persistent single-line footer with keybind hints (muted).

---

## Keybindings

| Key          | Action                                           |
|--------------|--------------------------------------------------|
| `j` / `down` | Move cursor down                                |
| `k` / `up`   | Move cursor up                                  |
| `enter`      | Copy selected item to system clipboard and exit  |
| `x`          | Clear entire clipboard history                   |
| `q` / `esc`  | Quit without changing clipboard                  |

---

## Architecture

```
prevy/
  main.go            -- entry point; initialises Bubble Tea program
  clipboard/
    clipboard.go     -- read/write system clipboard (xclip / xsel / wl-copy)
    history.go       -- in-memory history store + JSON persistence
  ui/
    model.go         -- Bubble Tea model (state, Init, Update, View)
    styles.go        -- Lip Gloss styles using Tokyo Night palette
    keys.go          -- key map definitions
  store/
    store.go         -- JSON file I/O (~/.local/share/prevy/history.json)
```

### Core Components

1. **Clipboard bridge** (`clipboard/clipboard.go`)
   - Detects available clipboard tool (`xclip`, `xsel`, `wl-copy`/`wl-paste`
     for Wayland, `pbcopy`/`pbpaste` for macOS).
   - `Read() (string, error)` -- reads current system clipboard.
   - `Write(text string) error` -- writes text to system clipboard.

2. **History store** (`clipboard/history.go` + `store/store.go`)
   - Keeps an ordered list of clipboard entries (content + timestamp).
   - Persists to `~/.local/share/prevy/history.json`.
   - `Add(entry)` -- prepends new entry, deduplicates.
   - `Clear()` -- wipes all entries and saves.
   - `Load() / Save()` -- JSON serialisation.

3. **TUI model** (`ui/model.go`)
   - Bubble Tea `Model` with fields: `entries`, `cursor`, `width`, `height`.
   - `Init()` -- loads history from disk, reads current clipboard, adds if new.
   - `Update()` -- handles key messages per the keymap.
   - `View()` -- renders the list with Lip Gloss styles.

4. **Styles** (`ui/styles.go`)
   - All Lip Gloss styles derived from the Tokyo Night palette table above.
   - Exported style variables: `TitleStyle`, `ItemStyle`, `SelectedStyle`,
     `IndexStyle`, `TimestampStyle`, `HelpStyle`, `BorderStyle`.

---

## Data Model

```go
type Entry struct {
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}
```

History is stored as a JSON array, newest first. Maximum 100 entries by default
(configurable later).

---

## User Flow

1. User copies text anywhere on the system.
2. User runs `prevy`.
3. Prevy reads the current clipboard; if the content is new, it prepends it to
   history.
4. The TUI renders the full history list.
5. User navigates with `j`/`k`, selects with `enter` (copies to clipboard and
   exits), presses `x` to clear all, or `q` to quit.

---

## Dependencies

| Package                                  | Purpose                  |
|------------------------------------------|--------------------------|
| `github.com/charmbracelet/bubbletea`     | TUI framework            |
| `github.com/charmbracelet/lipgloss`      | Styling / layout         |
| `github.com/atotto/clipboard`            | Cross-platform clipboard |

---

## MVP Scope (v0.1)

- [x] Masterplan
- [ ] Project scaffolding (`go mod init`, directory structure)
- [ ] Clipboard read/write bridge
- [ ] JSON-based history persistence
- [ ] Bubble Tea TUI with list, selection, clear, quit
- [ ] Tokyo Night styling with Lip Gloss
- [ ] `go build` produces single `prevy` binary

---

## Future Ideas (post-MVP)

- Background daemon that watches the clipboard and auto-records entries.
- Fuzzy search / filter (`/` key).
- Preview pane for long entries.
- Configurable max history size.
- Pin / favourite entries.
- Delete single entry (`d` key).
- Snap / AUR / Homebrew packaging.
