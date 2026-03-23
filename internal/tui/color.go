package tui

import (
	"os"
	"strings"
)

var forceColor bool

func SetForceColor(v bool) {
	forceColor = v
}

func SupportsColor() bool {
	if forceColor {
		return true
	}
	// Check --no-color flag via environment
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	return strings.Contains(term, "color")
}
