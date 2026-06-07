package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFilesGlob(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/a.txt", []byte("AAA"), 0644)
	_ = os.WriteFile(dir+"/b.txt", []byte("BBB"), 0644)

	results, err := ReadFiles(dir + "/*.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 files, got %d", len(results))
	}
}

func TestReadFilesMissing(t *testing.T) {
	_, err := ReadFiles("/nonexistent/*.go")
	if err == nil {
		t.Fatal("expected error for nonexistent pattern")
	}
}

func TestReadFilesRecursiveGlob(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	_ = os.MkdirAll(sub, 0755)
	_ = os.WriteFile(sub+"/c.txt", []byte("CCC"), 0644)

	results, err := ReadFiles(filepath.Join(dir, "**/*.txt"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 file via ** glob, got %d", len(results))
	}
}
