// Package compose generates Docker Compose files for CTF teams
package compose

import (
	"fmt"
	"path/filepath"

	"github.com/Lolozendev/CTFManager/internal/config"
	"github.com/Lolozendev/CTFManager/internal/model"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// Generator handles Docker Compose file generation
type Generator struct {
	config *config.Config
	logger *log.Logger
}

// New creates a new compose generator
func New(cfg *config.Config, logger *log.Logger) *Generator {
	return &Generator{
		config: cfg,
		logger: logger,
	}
}

// Generate creates a Docker Compose YAML for a team
func (g *Generator) Generate(team model.Team, challenges []model.Challenge) (string, error) {
	composeFile := model.NewComposeFile(team, challenges)

	// Marshal to YAML
	data, err := yaml.Marshal(&composeFile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal compose file: %w", err)
	}

	return string(data), nil
}

// PrepareTeamChallenges converts challenge list to team-specific configurations
func (g *Generator) PrepareTeamChallenges(challenges []model.Challenge, teamID int) []model.Challenge {
	prepared := make([]model.Challenge, len(challenges))

	for i, ch := range challenges {
		// Update paths to be absolute
		prepared[i] = model.Challenge{
			Name:      ch.Name,
			NetworkID: ch.NetworkID,
			BuildPath: g.config.GetChallengePath(
				model.FormatChallengeName(ch.NetworkID, ch.Name, ch.Enabled),
			),
			EnvPath: filepath.Join(
				g.config.GetChallengePath(
					model.FormatChallengeName(ch.NetworkID, ch.Name, ch.Enabled),
				),
				".env",
			),
			Enabled: ch.Enabled,
		}
	}

	return prepared
}
