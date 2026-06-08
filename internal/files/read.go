package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type FileResult struct {
	Path    string
	Content string
	Error   error
}

func ReadFiles(pattern string) ([]FileResult, error) {
	matches, err := doublestar.FilepathGlob(pattern)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}

	var results []FileResult
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil || info.IsDir() {
			continue
		}

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
		fmt.Fprintf(&sb, "<file path=\"%s\">%s</file>", r.Path, r.Content)
	}
	return sb.String()
}
