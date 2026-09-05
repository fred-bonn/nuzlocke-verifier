package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
)

func TestRunReturnsTheExpectedExitCodeForCLIArguments(t *testing.T) {
	tests := map[string]struct {
		args []string
		want int
	}{
		"rejects invalid weather":                 {args: []string{"-w", "99", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"rejects negative weather":                {args: []string{"-w", "-1", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"rejects non-positive iterations":         {args: []string{"--iterations", "0", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"rejects negative iterations":             {args: []string{"--iterations", "-5", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"rejects missing args":                    {args: []string{"--iterations", "1"}, want: 1},
		"rejects too many args":                   {args: []string{"data/player.txt", "data/rnb_trainer_1.txt", "extra.txt"}, want: 1},
		"rejects invalid showdown file":           {args: []string{"--iterations", "1", "data/nonexistent.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"rejects misconfigured policy file":       {args: []string{"--policy-file", "missing.json", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 1},
		"accepts valid CLI and returns zero":      {args: []string{"--iterations", "1", "data/player.txt", "data/rnb_trainer_1.txt"}, want: 0},
		"accepts a help request and returns zero": {args: []string{"-h"}, want: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := run(tc.args); got != tc.want {
				t.Fatalf("run(%v) = %d, want %d", tc.args, got, tc.want)
			}
		})
	}
}

func TestRunAcceptsVerboseFlag(t *testing.T) {
	oldVerbose := *verbose
	*verbose = false
	defer func() { *verbose = oldVerbose }()

	code := run([]string{"-v", "--iterations", "1", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 0 {
		t.Fatalf("expected zero exit code with verbose flag, got %d", code)
	}
	if !*verbose {
		t.Fatal("expected verbose flag to be enabled when parsing -v")
	}
}

func TestRunHelpDoesNotReturnError(t *testing.T) {
	if code := run([]string{"-h"}); code != 0 {
		t.Fatalf("expected zero exit code for help request, got %d", code)
	}
}

func TestParserEdgeCases(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"missing level line":     {input: "Horsea\nModest Nature\nAbility: Swift Swim\n- Bubble Beam\n", want: false},
		"invalid nature":         {input: "Horsea\nLevel: 17\nBad Nature\nAbility: Swift Swim\n- Bubble Beam\n", want: false},
		"invalid status":         {input: "Horsea\nLevel: 17\nModest Nature\nAbility: Swift Swim\nStatus> bad\n- Bubble Beam\n", want: false},
		"valid multi-move party": {input: "Horsea\nLevel: 17\nModest Nature\nAbility: Swift Swim\n- Bubble Beam\n- Twister\n", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "party.txt")
			if err := os.WriteFile(path, []byte(tc.input), 0o644); err != nil {
				t.Fatalf("write temp file: %v", err)
			}

			parsed, err := parser.ReadShowdownFile(path)
			if tc.want {
				if err != nil {
					t.Fatalf("expected valid parser result, got error: %v", err)
				}
				if len(parsed) == 0 {
					t.Fatal("expected at least one parsed pokemon")
				}
				if parsed[0].Name != "Horsea" {
					t.Fatalf("unexpected parsed name: %q", parsed[0].Name)
				}
				return
			}
			if err == nil {
				t.Fatal("expected parse error for malformed input")
			}
		})
	}
}

func TestRunRejectsInvalidPolicyFileFormat(t *testing.T) {
	dir := t.TempDir()
	badPolicy := filepath.Join(dir, "bad_policy.json")
	if err := os.WriteFile(badPolicy, []byte("{invalid json}"), 0o644); err != nil {
		t.Fatalf("write bad policy file: %v", err)
	}

	code := run([]string{"--policy-file", badPolicy, "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for invalid policy file, got %d", code)
	}
}
