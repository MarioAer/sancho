package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileResult struct {
	Path    string
	Content string
	Error   error
}

func ReadFiles(pattern string) ([]FileResult, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	if strings.Contains(pattern, "**") {
		dir := filepath.Dir(pattern)
		err := filepath.Walk(dir, func(path string, _ os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			ext := filepath.Ext(path)
			if ext == filepath.Ext(pattern) {
				matches = append(matches, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}

	seen := make(map[string]bool)
	var results []FileResult
	for _, m := range matches {
		if seen[m] {
			continue
		}
		seen[m] = true
		data, readErr := os.ReadFile(m)
		r := FileResult{Path: m}
		if readErr != nil {
			r.Error = readErr
		} else {
			r.Content = string(data)
		}
		results = append(results, r)
	}
	return results, nil
}

func FormatForPrompt(results []FileResult) string {
	var sb strings.Builder
	for _, r := range results {
		if r.Error != nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("<file path=\"%s\">%s</file>", r.Path, r.Content))
	}
	return sb.String()
}
