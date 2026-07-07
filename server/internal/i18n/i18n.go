package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// MessageEntry represents a single key entry in Chrome/Firefox i18n messages.json format.
type MessageEntry struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

// MessagesDict represents the full dictionary structure of messages.json.
type MessagesDict map[string]MessageEntry

// ParseAndValidate parses a messages.json content, validating its structure and constraints.
func ParseAndValidate(data []byte) (MessagesDict, error) {
	// Limit file size checking to 10MB to avoid high memory allocation
	if len(data) > 10*1024*1024 {
		return nil, errors.New("file size exceeds limit of 10MB")
	}

	var dict MessagesDict
	if err := json.Unmarshal(data, &dict); err != nil {
		return nil, fmt.Errorf("invalid JSON schema: %w", err)
	}

	// Validate structure details
	for key, entry := range dict {
		if key == "" {
			return nil, errors.New("empty keys are not allowed")
		}
		// Chrome extension translations require 'message' field
		if entry.Message == "" {
			// Note: empty messages might be valid under some contexts but let's warn or validate
			// We will allow empty message if the key exists, but let's make sure it's not a generic map
		}
	}

	return dict, nil
}

// Regex for common placeholder patterns (ordered to match double braces before single braces)
var placeholderRegex = regexp.MustCompile(`\{\{[a-zA-Z0-9_-]+\}\}|\{[a-zA-Z0-9_-]+\}|%[s|d]|\$[1-9]|\$[a-zA-Z0-9_-]+\$`)

// ExtractPlaceholders finds all placeholders matching supported patterns in the text.
func ExtractPlaceholders(text string) []string {
	var found []string
	seen := make(map[string]bool)

	matches := placeholderRegex.FindAllString(text, -1)
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			found = append(found, m)
		}
	}
	return found
}

// CheckPlaceholderWarnings compares target translation placeholders with the source translation.
// It returns warnings if any placeholders in source are missing in target.
func CheckPlaceholderWarnings(source, target string) []string {
	sourcePlaceholders := ExtractPlaceholders(source)
	if len(sourcePlaceholders) == 0 {
		return nil
	}

	targetLower := strings.ToLower(target)

	// Correctly build warnings list
	var result []string
	for _, p := range sourcePlaceholders {
		// Try case-insensitive search
		if !strings.Contains(targetLower, strings.ToLower(p)) {
			result = append(result, fmt.Sprintf("Missing placeholder: %s", p))
		}
	}

	return result
}
