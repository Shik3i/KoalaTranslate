package i18n

import (
	"testing"
)

func TestParseAndValidate(t *testing.T) {
	validJSON := []byte(`{
		"appName": {
			"message": "KoalaSync",
			"description": "App Name"
		},
		"emptyDesc": {
			"message": "Hello"
		}
	}`)

	dict, err := ParseAndValidate(validJSON)
	if err != nil {
		t.Fatalf("expected valid JSON to parse without error, got: %v", err)
	}

	if len(dict) != 2 {
		t.Errorf("expected 2 entries, got %d", len(dict))
	}

	if dict["appName"].Message != "KoalaSync" {
		t.Errorf("expected Message 'KoalaSync', got '%s'", dict["appName"].Message)
	}

	invalidJSON := []byte(`{
		"appName": {
			"message":
		}
	}`)

	_, err = ParseAndValidate(invalidJSON)
	if err == nil {
		t.Error("expected invalid JSON to return an error, got nil")
	}

	largeJSON := make([]byte, 11*1024*1024) // 11MB
	_, err = ParseAndValidate(largeJSON)
	if err == nil || err.Error() != "file size exceeds limit of 10MB" {
		t.Errorf("expected file size limit error, got: %v", err)
	}
}

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hello {name}!", []string{"{name}"}},
		{"Hello {{name}} and {{age}}", []string{"{{name}}", "{{age}}"}},
		{"Count: %d, string: %s", []string{"%d", "%s"}},
		{"Result $1 and $2", []string{"$1", "$2"}},
		{"Variables $VAL$ and $VAL2$", []string{"$VAL$", "$VAL2$"}},
		{"No placeholders here", nil},
	}

	for _, tc := range tests {
		actual := ExtractPlaceholders(tc.input)
		if len(actual) != len(tc.expected) {
			t.Errorf("for input %q: expected %v, got %v", tc.input, tc.expected, actual)
			continue
		}
		for i, v := range actual {
			if v != tc.expected[i] {
				t.Errorf("for input %q: expected %v, got %v", tc.input, tc.expected, actual)
			}
		}
	}
}

func TestCheckPlaceholderWarnings(t *testing.T) {
	tests := []struct {
		source   string
		target   string
		expected int // number of warnings
	}{
		{"Hello {name}!", "Hallo {name}!", 0},
		{"Hello {name}!", "Hallo name!", 1},
		{"Wait {seconds}s", "Warten Sie {seconds} Sek", 0},
		{"Wait {seconds}s", "Warten Sie Sekunden", 1},
		{"Match $1 and $2", "Match $1 und $2", 0},
		{"Match $1 and $2", "Match $1", 1},
	}

	for _, tc := range tests {
		warnings := CheckPlaceholderWarnings(tc.source, tc.target)
		if len(warnings) != tc.expected {
			t.Errorf("for source %q and target %q: expected %d warnings, got %d (%v)",
				tc.source, tc.target, tc.expected, len(warnings), warnings)
		}
	}
}
