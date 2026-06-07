package cmd

import (
	"testing"
)

func TestRootCommandBuild(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("expected non-nil root command")
	}
	if cmd.Use != "sancho" {
		t.Fatalf("expected use=sancho, got %s", cmd.Use)
	}
}

func TestRootSubcommands(t *testing.T) {
	root := NewRootCmd()

	want := map[string]bool{
		"ask":   false,
		"write": false,
	}
	for _, c := range root.Commands() {
		if _, ok := want[c.Use]; ok {
			want[c.Use] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("missing subcommand: %s", name)
		}
	}
}
