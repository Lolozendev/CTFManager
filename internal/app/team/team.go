// Package team provides functionality for managing CTF teams
package team

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Lolozendev/CTFManager/internal/config"
	"github.com/Lolozendev/CTFManager/internal/model"
	"github.com/charmbracelet/log"
)

// Manager handles team operations
type Manager struct {
	config *config.Config
	logger *log.Logger
}

// New creates a new team manager
func New(cfg *config.Config, logger *log.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// Create creates a new team
func (m *Manager) Create(id int, name string, members []string) error {
	// Validate team ID
	if id < m.config.Teams.MinID || id > m.config.Teams.MaxID {
		return fmt.Errorf("invalid team ID %d (must be between %d and %d)",
			id, m.config.Teams.MinID, m.config.Teams.MaxID)
	}

	// Check if team already exists
	teamPath := m.config.GetTeamPath(fmt.Sprintf("%d-%s", id, name))
	if _, err := os.Stat(teamPath); !os.IsNotExist(err) {
		return fmt.Errorf("team %s already exists", name)
	}

	// Create team directory
	if err := os.MkdirAll(teamPath, 0755); err != nil {
		return fmt.Errorf("failed to create team directory: %w", err)
	}

	m.logger.Info("Team created", "id", id, "name", name, "members", len(members))
	return nil
}

// List returns all teams
func (m *Manager) List() ([]model.Team, error) {
	entries, err := os.ReadDir(m.config.Paths.Teams)
	if err != nil {
		return nil, fmt.Errorf("failed to read teams directory: %w", err)
	}

	var teams []model.Team
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		id, name, enabled, err := model.ParseChallengeName(entry.Name()) // Reuse same parsing logic
		if err != nil {
			m.logger.Warn("Skipping invalid team directory", "name", entry.Name(), "error", err)
			continue
		}

		team := model.Team{
			ID:      id,
			Name:    name,
			Enabled: enabled,
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// Delete removes a team
func (m *Manager) Delete(name string) error {
	teams, err := m.List()
	if err != nil {
		return err
	}

	var found *model.Team
	for _, t := range teams {
		if t.Name == name && t.Enabled {
			found = &t
			break
		}
	}

	if found == nil {
		return fmt.Errorf("team %s not found", name)
	}

	teamPath := m.config.GetTeamPath(model.FormatChallengeName(found.ID, name, true))
	if err := os.RemoveAll(teamPath); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	m.logger.Info("Team deleted", "name", name)
	return nil
}

// Disable disables a team
func (m *Manager) Disable(name string) error {
	teams, err := m.List()
	if err != nil {
		return err
	}

	var found *model.Team
	for _, t := range teams {
		if t.Name == name && t.Enabled {
			found = &t
			break
		}
	}

	if found == nil {
		return fmt.Errorf("enabled team %s not found", name)
	}

	oldPath := m.config.GetTeamPath(model.FormatChallengeName(found.ID, name, true))
	newPath := m.config.GetTeamPath(model.FormatChallengeName(0, name, false))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to disable team: %w", err)
	}

	m.logger.Info("Team disabled", "name", name)
	return nil
}

// Enable enables a team
func (m *Manager) Enable(name string, id int) error {
	oldPath := m.config.GetTeamPath(fmt.Sprintf("x-%s", name))
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("disabled team %s not found", name)
	}

	newPath := m.config.GetTeamPath(model.FormatChallengeName(id, name, true))
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to enable team: %w", err)
	}

	m.logger.Info("Team enabled", "name", name, "id", id)
	return nil
}

// Validate checks if a team directory has the required structure
func (m *Manager) Validate(teamName string) error {
	teams, err := m.List()
	if err != nil {
		return err
	}

	var found *model.Team
	for _, t := range teams {
		if t.Name == teamName && t.Enabled {
			found = &t
			break
		}
	}

	if found == nil {
		return fmt.Errorf("team %s not found", teamName)
	}

	teamPath := m.config.GetTeamPath(model.FormatChallengeName(found.ID, teamName, true))
	composePath := filepath.Join(teamPath, "compose.yml")

	if _, err := os.Stat(composePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("team %s missing compose.yml file", teamName)
		}
		return fmt.Errorf("error checking compose.yml: %w", err)
	}

	return nil
}
