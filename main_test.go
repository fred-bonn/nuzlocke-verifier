package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
)

func TestRunRejectsInvalidWeather(t *testing.T) {
	code := run([]string{"-w", "99", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for invalid weather, got %d", code)
	}
}

func TestRunRejectsNonPositiveIterations(t *testing.T) {
	code := run([]string{"--iterations", "0", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for invalid iteration count, got %d", code)
	}
}

func TestRunRejectsMissingArgs(t *testing.T) {
	code := run([]string{"--iterations", "1"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for missing arguments, got %d", code)
	}
}

func TestRunAcceptsValidCLIAndReturnsZero(t *testing.T) {
	code := run([]string{"--iterations", "1", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 0 {
		t.Fatalf("expected zero exit code for valid input, got %d", code)
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

func TestRunRejectsMisconfiguredPolicyFile(t *testing.T) {
	code := run([]string{"--policy-file", "missing.json", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for missing policy file, got %d", code)
	}
}

func TestParserEdgeCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "missing level line", input: "Horsea\nModest Nature\nAbility: Swift Swim\n- Bubble Beam\n", want: false},
		{name: "invalid nature", input: "Horsea\nLevel: 17\nBad Nature\nAbility: Swift Swim\n- Bubble Beam\n", want: false},
		{name: "invalid status", input: "Horsea\nLevel: 17\nModest Nature\nAbility: Swift Swim\nStatus> bad\n- Bubble Beam\n", want: false},
		{name: "valid multi-move party", input: "Horsea\nLevel: 17\nModest Nature\nAbility: Swift Swim\n- Bubble Beam\n- Twister\n", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestRunRejectsNegativeWeather(t *testing.T) {
	code := run([]string{"-w", "-1", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for negative weather, got %d", code)
	}
}

func TestRunRejectsNegativeIterations(t *testing.T) {
	code := run([]string{"--iterations", "-5", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for negative iterations, got %d", code)
	}
}

func TestRunRejectsInvalidShowdownFile(t *testing.T) {
	code := run([]string{"--iterations", "1", "data/nonexistent.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for nonexistent file, got %d", code)
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

func TestRunRejectsZeroIterations(t *testing.T) {
	code := run([]string{"--iterations", "0", "data/player.txt", "data/rnb_trainer_1.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for zero iterations, got %d", code)
	}
}

func TestRunRejectsTooManyArgs(t *testing.T) {
	code := run([]string{"data/player.txt", "data/rnb_trainer_1.txt", "extra.txt"})
	if code != 1 {
		t.Fatalf("expected exit code 1 for too many arguments, got %d", code)
	}
}
