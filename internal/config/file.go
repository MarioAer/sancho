package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

func LoadFile(dir string) (Config, error) {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}

	candidates := []string{
		filepath.Join(dir, ".sancho.json"),
		filepath.Join(home, ".config", "sancho", "config.json"),
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		cleaned := stripJSONC(string(data))
		var cfg Config
		if err := json.Unmarshal([]byte(cleaned), &cfg); err != nil {
			return Config{}, err
		}
		return cfg, nil
	}

	return Config{}, nil
}

var lineComment = regexp.MustCompile(`(?m)//[^\n]*`)
var blockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)

func stripJSONC(s string) string {
	s = lineComment.ReplaceAllString(s, "")
	s = blockComment.ReplaceAllString(s, "")
	return s
}
