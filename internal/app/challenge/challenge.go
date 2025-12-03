// Package challenge provides functionality for managing CTF challenges
package challenge

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Lolozendev/CTFManager/internal/config"
	"github.com/Lolozendev/CTFManager/internal/model"
	"github.com/charmbracelet/log"
)

var (
	challengeNameRegexp = regexp.MustCompile(`^(?:(\d{1,3})|x)-(\w+)$`)
)

// Manager handles challenge operations
type Manager struct {
	config *config.Config
	logger *log.Logger
}

// New creates a new challenge manager
func New(cfg *config.Config, logger *log.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// List returns all challenges found in the challenges directory
func (m *Manager) List() ([]model.Challenge, error) {
	entries, err := os.ReadDir(m.config.Paths.Challenges)
	if err != nil {
		return nil, fmt.Errorf("failed to read challenges directory: %w", err)
	}

	var challenges []model.Challenge
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		networkID, name, enabled, err := model.ParseChallengeName(entry.Name())
		if err != nil {
			m.logger.Warn("Skipping invalid challenge directory", "name", entry.Name(), "error", err)
			continue
		}

		challengePath := m.config.GetChallengePath(entry.Name())

		challenge := model.Challenge{
			Name:      name,
			NetworkID: networkID,
			BuildPath: challengePath,
			EnvPath:   filepath.Join(challengePath, ".env"),
			Enabled:   enabled,
		}

		challenges = append(challenges, challenge)
	}

	return challenges, nil
}

// ListEnabled returns only enabled challenges
func (m *Manager) ListEnabled() ([]model.Challenge, error) {
	allChallenges, err := m.List()
	if err != nil {
		return nil, err
	}

	var enabled []model.Challenge
	for _, ch := range allChallenges {
		if ch.Enabled {
			enabled = append(enabled, ch)
		}
	}

	return enabled, nil
}

// Validate checks all challenges for correctness
func (m *Manager) Validate() error {
	challenges, err := m.ListEnabled()
	if err != nil {
		return err
	}

	if len(challenges) == 0 {
		return errors.New("no enabled challenges found")
	}

	// Check for duplicate network IDs
	usedIDs := make(map[int]string)
	for _, ch := range challenges {
		if ch.NetworkID < m.config.Challenges.MinNetworkID ||
			ch.NetworkID > m.config.Challenges.MaxNetworkID {
			return fmt.Errorf("challenge %s has invalid network ID %d (must be between %d and %d)",
				ch.Name, ch.NetworkID,
				m.config.Challenges.MinNetworkID,
				m.config.Challenges.MaxNetworkID)
		}

		if existingName, exists := usedIDs[ch.NetworkID]; exists {
			return fmt.Errorf("duplicate network ID %d used by challenges: %s and %s",
				ch.NetworkID, existingName, ch.Name)
		}
		usedIDs[ch.NetworkID] = ch.Name

		// Validate directory structure
		if err := m.validateChallengeStructure(ch); err != nil {
			return fmt.Errorf("challenge %s: %w", ch.Name, err)
		}
	}

	return nil
}

// validateChallengeStructure checks if a challenge has required files
func (m *Manager) validateChallengeStructure(ch model.Challenge) error {
	// Check for Dockerfile
	dockerfilePath := filepath.Join(ch.BuildPath, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("missing Dockerfile")
		}
		return fmt.Errorf("error checking Dockerfile: %w", err)
	}

	// Check for .env file
	if _, err := os.Stat(ch.EnvPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("missing .env file (create an empty one if not needed)")
		}
		return fmt.Errorf("error checking .env file: %w", err)
	}

	return nil
}

// Enable enables a challenge by renaming its directory
func (m *Manager) Enable(name string, networkID int) error {
	// Find the disabled challenge
	oldPath := m.config.GetChallengePath(fmt.Sprintf("x-%s", name))
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("disabled challenge %s not found", name)
	}

	// Check if network ID is available
	challenges, err := m.ListEnabled()
	if err != nil {
		return err
	}

	for _, ch := range challenges {
		if ch.NetworkID == networkID {
			return fmt.Errorf("network ID %d is already used by challenge %s", networkID, ch.Name)
		}
	}

	// Rename directory
	newPath := m.config.GetChallengePath(model.FormatChallengeName(networkID, name, true))
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to enable challenge: %w", err)
	}

	m.logger.Info("Challenge enabled", "name", name, "networkID", networkID)
	return nil
}

// Disable disables a challenge by renaming its directory
func (m *Manager) Disable(name string) error {
	// Find the enabled challenge
	challenges, err := m.ListEnabled()
	if err != nil {
		return err
	}

	var found *model.Challenge
	for _, ch := range challenges {
		if ch.Name == name {
			found = &ch
			break
		}
	}

	if found == nil {
		return fmt.Errorf("enabled challenge %s not found", name)
	}

	// Rename directory
	oldPath := found.BuildPath
	newPath := m.config.GetChallengePath(model.FormatChallengeName(0, name, false))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to disable challenge: %w", err)
	}

	m.logger.Info("Challenge disabled", "name", name)
	return nil
}
