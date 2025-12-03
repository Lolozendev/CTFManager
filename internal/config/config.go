package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all configuration for CTFManager
type Config struct {
	Paths      PathConfig
	Network    NetworkConfig
	Challenges ChallengeConfig
	Teams      TeamConfig
}

// PathConfig defines file system paths
type PathConfig struct {
	Challenges      string
	Teams           string
	DnsmasqTemplate string
}

// NetworkConfig defines network ranges
type NetworkConfig struct {
	BaseSubnet string // e.g., "10.0"
}

// ChallengeConfig defines challenge constraints
type ChallengeConfig struct {
	MinNetworkID int
	MaxNetworkID int
}

// TeamConfig defines team constraints
type TeamConfig struct {
	MinID       int
	MaxID       int
	BaseVPNPort int // Base port for VPN (e.g., 50000)
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Paths: PathConfig{
			Challenges:      "/challenges",
			Teams:           "/equipes",
			DnsmasqTemplate: "/dnsconf/dnsmasq.template",
		},
		Network: NetworkConfig{
			BaseSubnet: "10.0",
		},
		Challenges: ChallengeConfig{
			MinNetworkID: 11,
			MaxNetworkID: 249,
		},
		Teams: TeamConfig{
			MinID:       1,
			MaxID:       254,
			BaseVPNPort: 50000,
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check if challenge path exists
	if _, err := os.Stat(c.Paths.Challenges); err != nil {
		return fmt.Errorf("challenge path %s does not exist: %w", c.Paths.Challenges, err)
	}

	// Check if teams path exists (create if not)
	if _, err := os.Stat(c.Paths.Teams); os.IsNotExist(err) {
		if err := os.MkdirAll(c.Paths.Teams, 0755); err != nil {
			return fmt.Errorf("failed to create teams directory: %w", err)
		}
	}

	// Validate network ranges
	if c.Challenges.MinNetworkID >= c.Challenges.MaxNetworkID {
		return fmt.Errorf("invalid challenge network range: min=%d max=%d",
			c.Challenges.MinNetworkID, c.Challenges.MaxNetworkID)
	}

	if c.Teams.MinID >= c.Teams.MaxID {
		return fmt.Errorf("invalid team ID range: min=%d max=%d",
			c.Teams.MinID, c.Teams.MaxID)
	}

	return nil
}

// GetChallengePath returns the full path to a challenge directory
func (c *Config) GetChallengePath(challengeName string) string {
	return filepath.Join(c.Paths.Challenges, challengeName)
}

// GetTeamPath returns the full path to a team directory
func (c *Config) GetTeamPath(teamName string) string {
	return filepath.Join(c.Paths.Teams, teamName)
}

// GetVPNPort returns the VPN port for a team
func (c *Config) GetVPNPort(teamID int) int {
	return c.Teams.BaseVPNPort + teamID
}
