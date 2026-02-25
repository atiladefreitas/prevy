# AGENTS.md - Prevy

## Project Overview

Prevy is a minimal terminal clipboard history manager written in Go using
[Charm Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI framework)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss) (styling). It has two
modes: a background daemon (`prevy --daemon`) that polls the system clipboard and
persists entries to `~/.local/share/prevy/history.json`, and a TUI (`prevy`) that
displays the clipboard history and lets users copy, paste, or clear entries.

Module path: `github.com/atiladefreitas/prevy`

## Directory Structure

```
prevy/
  main.go              Entry point: CLI arg parsing, Bubble Tea program launch
  clipboard/
    clipboard.go       System clipboard read/write (platform detection via os/exec)
  daemon/
    daemon.go          Background clipboard watcher: PID management, poll loop, signals
  store/
    store.go           JSON file persistence: Load, Save, Add (deduplicate + cap), Clear
  ui/
    model.go           Bubble Tea Model: Init, Update, View; rendering, scrolling
    keys.go            Keymap: parseKey mapping key strings to keyAction enum
    styles.go          Lip Gloss styles: Tokyo Night color constants + exported style vars
```

## Build / Run / Test Commands

There is no Makefile. Use standard Go tooling:

```bash
# Build
go build -o prevy .

# Run TUI
./prevy

# Run daemon
./prevy --daemon

# Install globally
go install github.com/atiladefreitas/prevy@latest

# Run all tests
go test ./...

# Run tests in a single package
go test ./store
go test ./clipboard
go test ./daemon
go test ./ui

# Run a single test by name
go test ./store -run TestAdd

# Verbose test output
go test -v ./...

# Vet (static analysis)
go vet ./...

# Format check (gofmt uses tabs; no custom config exists)
gofmt -l .
```

No CI/CD pipeline, linter config, or pre-commit hooks exist. Code should pass
`go vet ./...` and `gofmt -l .` with no output. No tests exist yet -- when
adding features, add corresponding `*_test.go` files in the relevant package.

## Code Style Guidelines

### Formatting

Standard `gofmt` formatting. Indentation uses tabs (Go default). No `.editorconfig`
or custom formatter configuration. Line length is not enforced but kept reasonable
(~80-100 chars, with Lip Gloss chains being an exception).

### Import Organization

Three groups separated by blank lines:

1. Standard library
2. Third-party packages
3. Internal project packages (`github.com/atiladefreitas/prevy/...`)

```go
import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/atiladefreitas/prevy/daemon"
    "github.com/atiladefreitas/prevy/ui"
)
```

The `bubbletea` package is aliased as `tea` (Charm community convention). Import it
as `tea "github.com/charmbracelet/bubbletea"` or just `"github.com/charmbracelet/bubbletea"`
(Go resolves the package name to `tea` either way). Prefer the explicit alias in
files that use it heavily.

### Naming Conventions

- **Packages**: lowercase, single-word (`clipboard`, `daemon`, `store`, `ui`)
- **Exported symbols**: PascalCase (`Read`, `Write`, `Load`, `Save`, `Model`, `New`)
- **Unexported symbols**: camelCase (`detect`, `dataPath`, `parseKey`, `relativeTime`)
- **Constants**: camelCase for unexported (`maxEntries`, `pollInterval`); iota enums
  use camelCase (`statusBrowsing`, `keyUp`, `keyQuit`)
- **Unexported color vars**: camelCase (`background`, `surface`, `blue`, `cyan`)
- **Exported style vars**: PascalCase (`AppStyle`, `TitleStyle`, `BorderStyle`)
- **Struct fields**: camelCase for unexported (`entries`, `cursor`, `width`)
- **JSON tags**: lowercase (`json:"content"`, `json:"timestamp"`)

### Type Patterns

- Enum-like types: custom `int` types with `iota` constants (`status`, `keyAction`)
- Data structs with JSON tags for persistence (`store.Entry`)
- The TUI model (`ui.Model`) satisfies `tea.Model` implicitly (no explicit interface)
- Use named fields in struct literals, multi-line

### Error Handling

The codebase has two patterns depending on context:

- **In the UI layer**: errors from clipboard/store operations are silenced with `_`
  because the TUI must not crash:
  ```go
  _ = clipboard.Write(m.entries[m.cursor].Content)
  _ = store.Clear()
  ```
- **In daemon/store**: errors are checked and returned early with context via
  `fmt.Errorf("...: %w", err)`:
  ```go
  if err != nil {
      return fmt.Errorf("failed to write pid file: %w", err)
  }
  ```
- Graceful fallback with `os.IsNotExist` when history file does not exist yet
- Fatal errors in `main.go` use `fmt.Fprintf(os.Stderr, ...)` + `os.Exit(1)`
- The clipboard bridge returns `("", nil)` when no provider is found (no error)

When adding new code: prefer returning errors over silencing them. Only silence
errors in the Bubble Tea Update/View methods where returning an error is not
possible. Always wrap errors with `fmt.Errorf("context: %w", err)`.

### File Permissions

Use Go octal literals with the `0o` prefix: `0o755` for directories, `0o644` for
files.

### Lip Gloss Style Conventions

- All colors use the Tokyo Night palette defined in `ui/styles.go`
- Colors are package-level `lipgloss.Color` variables (unexported)
- Styles are package-level `var` declarations using the fluent builder pattern:
  ```go
  var TitleStyle = lipgloss.NewStyle().
      Foreground(blue).
      Bold(true)
  ```
- Do not hardcode hex values outside `ui/styles.go` -- use the named color variables

### Bubble Tea Architecture (Elm Architecture)

The TUI follows Bubble Tea's Elm-style pattern:
- `Model` struct holds all state (no global state)
- `Init() tea.Cmd` returns initial command (currently nil)
- `Update(msg tea.Msg) (tea.Model, tea.Cmd)` handles messages via type switch
- `View() string` renders the entire UI as a string
- Key handling is decoupled through `parseKey()` in `keys.go`
- No dependency injection; packages call each other's exported functions directly

### Platform Support

The clipboard bridge (`clipboard/clipboard.go`) supports:
- **macOS**: `pbcopy`/`pbpaste`
- **Linux Wayland**: `wl-copy`/`wl-paste` (preferred)
- **Linux X11**: `xclip` or `xsel` (fallback)

Detection is done at runtime via `exec.LookPath`. When adding clipboard
functionality, maintain this provider pattern and priority order.

### Data Storage

- History file: `~/.local/share/prevy/history.json`
- PID file: `~/.local/share/prevy/daemon.pid`
- Maximum 100 entries (constant `maxEntries` in `store/store.go`)
- Entries are stored as JSON array, newest first
- Deduplication on add: existing entry with same content is removed before prepend
