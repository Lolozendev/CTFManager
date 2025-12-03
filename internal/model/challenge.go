package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Challenge represents a CTF challenge
type Challenge struct {
	Name      string
	NetworkID int    // Network position (11-249)
	BuildPath string // Path to Dockerfile
	EnvPath   string // Path to .env file
	Enabled   bool
}

// ParseChallengeName parses a challenge directory name (format: "11-webchallenge" or "x-disabled")
// Returns: networkID, name, enabled, error
func ParseChallengeName(dirName string) (int, string, bool, error) {
	parts := strings.SplitN(dirName, "-", 2)
	if len(parts) != 2 {
		return 0, "", false, fmt.Errorf("invalid challenge name format: %s (expected: <number>-<name> or x-<name>)", dirName)
	}

	prefix := parts[0]
	name := parts[1]

	// Check if disabled
	if prefix == "x" {
		return 0, name, false, nil
	}

	// Parse network ID
	networkID, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, "", false, fmt.Errorf("invalid challenge number in %s: %w", dirName, err)
	}

	return networkID, name, true, nil
}

// FormatChallengeName creates a directory name from challenge data
func FormatChallengeName(networkID int, name string, enabled bool) string {
	if !enabled {
		return fmt.Sprintf("x-%s", name)
	}
	return fmt.Sprintf("%d-%s", networkID, name)
}
