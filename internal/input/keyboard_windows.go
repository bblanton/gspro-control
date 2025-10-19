//go:build windows

package input

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCombo presses keys in combo like ["ctrl", "m"].
// Windows implementation using PowerShell SendKeys via .NET Windows Forms.
// It targets the currently focused window (ensure GSPro has focus).
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

	sendKeys := buildSendKeysString(mods, key)

	ps := buildPowerShellCommand(sendKeys)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("powershell SendKeys failed: %v: %s", err, stderr.String())
	}
	return nil
}

// buildSendKeysString converts modifiers+key into a Windows SendKeys string.
// Modifiers: ctrl(^), alt(%), shift(+). The Windows key is not supported by SendKeys and is ignored.
func buildSendKeysString(mods []string, key string) string {
	var b strings.Builder
	for _, m := range mods {
		switch m {
		case "ctrl", "control":
			b.WriteString("^")
		case "alt", "menu", "option":
			b.WriteString("%")
		case "shift":
			b.WriteString("+")
		case "cmd", "win", "windows", "meta":
			// SendKeys doesn't support the Windows key; ignore
		}
	}

	b.WriteString(escapeSendKeysToken(key))
	return b.String()
}

// escapeSendKeysToken escapes a single key token for SendKeys.
// For named keys it returns the appropriate token like {TAB}, {ESC}, {LEFT}.
// For single characters, it wraps special characters in braces per SendKeys rules.
func escapeSendKeysToken(k string) string {
	switch k {
	case "enter", "return":
		return "~" // Enter in SendKeys
	case "tab":
		return "{TAB}"
	case "escape", "esc":
		return "{ESC}"
	case "space", "spacebar":
		return "{SPACE}"
	case "left":
		return "{LEFT}"
	case "right":
		return "{RIGHT}"
	case "up":
		return "{UP}"
	case "down":
		return "{DOWN}"
	}

	if len(k) == 1 {
		ch := k
		// Special characters for SendKeys that must be wrapped in braces
		if strings.ContainsAny(ch, "+^%~(){}") {
			return "{" + ch + "}"
		}
		return ch
	}

	// Default to braces for unknown named keys, e.g., {F5}
	upper := strings.ToUpper(k)
	return "{" + upper + "}"
}

// buildPowerShellCommand returns a PS command that loads Windows Forms and sends keys.
// We use single quotes around the SendKeys string and escape any single quotes by doubling.
func buildPowerShellCommand(sendKeys string) string {
	// Escape single quotes for single-quoted PS string
	escaped := strings.ReplaceAll(sendKeys, "'", "''")
	// Construct a robust command with error stop
	// Add-Type may already be loaded; that's fine.
	return fmt.Sprintf("$ErrorActionPreference='Stop'; Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait('%s')", escaped)
}
