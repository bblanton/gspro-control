//go:build darwin

package input

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCombo presses keys in combo like ["ctrl", "m"].
// macOS implementation using AppleScript via osascript.
func ExecuteCombo(combo []string) error {
	if len(combo) == 0 {
		return errors.New("empty combo")
	}

	normalized := make([]string, 0, len(combo))
	for _, k := range combo {
		normalized = append(normalized, strings.ToLower(strings.TrimSpace(k)))
	}

	key := normalized[len(normalized)-1]
	mods := normalized[:len(normalized)-1]

	// Map modifier synonyms to AppleScript names
	var usingParts []string
	for _, m := range mods {
		switch m {
		case "ctrl", "control":
			usingParts = append(usingParts, "control down")
		case "alt", "option":
			usingParts = append(usingParts, "option down")
		case "shift":
			usingParts = append(usingParts, "shift down")
		case "cmd", "command", "meta":
			usingParts = append(usingParts, "command down")
		default:
			// ignore unknown modifiers
		}
	}

	// Build AppleScript
	var script string
	if len(usingParts) > 0 {
		script = fmt.Sprintf("tell application \"System Events\" to keystroke %s using {%s}", appleScriptKeyLiteral(key), strings.Join(usingParts, ", "))
	} else {
		script = fmt.Sprintf("tell application \"System Events\" to keystroke %s", appleScriptKeyLiteral(key))
	}

	cmd := exec.Command("osascript", "-e", script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("osascript failed: %v: %s", err, stderr.String())
	}
	return nil
}

// appleScriptKeyLiteral returns an AppleScript-safe literal for keystroke.
// For single-character keys, it returns a quoted string, e.g., "\"m\"".
// For some named keys, it maps to key code when necessary.
func appleScriptKeyLiteral(k string) string {
	// Common named keys mapping if needed in future.
	// For now, pass most characters directly as a string.
	if len(k) == 1 {
		// Escape backslash and quotes
		escaped := strings.ReplaceAll(k, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", escaped)
	}

	// Basic named keys support
	switch k {
	case "enter", "return":
		return "return"
	case "tab":
		return "tab"
	case "escape", "esc":
		return "escape"
	case "space", "spacebar":
		return "space"
	case "left", "right", "up", "down":
		// Arrow keys require key code on macOS
		// left=123, right=124, down=125, up=126
		code := map[string]int{"left": 123, "right": 124, "down": 125, "up": 126}[k]
		return fmt.Sprintf("key code %d", code)
	}

	// Default: quote it
	escaped := strings.ReplaceAll(k, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", escaped)
}
