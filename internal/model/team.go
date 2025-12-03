package model

// Member represents a team member
type Member struct {
	Username string `json:"username"`
}

// Team represents a CTF team with its infrastructure
type Team struct {
	ID      int
	Name    string
	Members []Member
	Enabled bool
}

// ComposeFile represents a complete Docker Compose configuration
type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks"`
}

// NewComposeFile creates a Docker Compose configuration for a team
func NewComposeFile(team Team, challenges []Challenge) ComposeFile {
	networkName := team.Name + "-Network"

	services := make(map[string]Service)

	// Add infrastructure services
	services["wireguard"] = NewWireguardService(team.Name, team.ID, len(team.Members))
	services["dnsmasq"] = NewDnsmasqService(team.Name, team.ID)

	// Add challenge services
	for _, challenge := range challenges {
		services[challenge.Name] = NewChallengeService(
			team.Name,
			team.ID,
			challenge.NetworkID,
			challenge.Name,
			challenge.BuildPath,
			challenge.EnvPath,
		)
	}

	networks := make(map[string]Network)
	networks[networkName] = NewTeamNetwork(team.ID)

	return ComposeFile{
		Services: services,
		Networks: networks,
	}
}
