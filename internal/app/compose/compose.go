// Package compose generates Docker Compose files for CTF teams
package compose

import (
	"fmt"

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
