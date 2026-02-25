package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

const maxEntries = 100

func dataPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, ".local", "share", "prevy")
	if err := os.MkdirAll(p, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(p, "history.json"), nil
}

func Load() ([]Entry, error) {
	p, err := dataPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return []Entry{}, nil
	}
	return entries, nil
}

func Save(entries []Entry) error {
	p, err := dataPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func Add(entries []Entry, content string) []Entry {
	content = trimContent(content)
	if content == "" {
		return entries
	}

	// deduplicate: remove existing entry with same content
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if e.Content != content {
			filtered = append(filtered, e)
		}
	}

	entry := Entry{
		Content:   content,
		Timestamp: time.Now(),
	}
	result := append([]Entry{entry}, filtered...)

	if len(result) > maxEntries {
		result = result[:maxEntries]
	}
	return result
}

func Clear() error {
	return Save([]Entry{})
}

func trimContent(s string) string {
	// trim leading/trailing whitespace but preserve internal structure
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
