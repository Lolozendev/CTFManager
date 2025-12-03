package model

// Network represents Docker Compose network configuration
type Network struct {
	Driver string      `yaml:"driver"`
	IPAM   NetworkIPAM `yaml:"ipam"`
}

// NetworkIPAM represents IP Address Management configuration
type NetworkIPAM struct {
	Config []NetworkConfig `yaml:"config"`
}

// NetworkConfig represents network subnet configuration
type NetworkConfig struct {
	Subnet  string `yaml:"subnet"`
	Gateway string `yaml:"gateway"`
}

// NewTeamNetwork creates a network for a team
func NewTeamNetwork(teamNumber int) Network {
	return Network{
		Driver: "bridge",
		IPAM: NetworkIPAM{
			Config: []NetworkConfig{
				{
					Subnet:  formatStr("10.0.%d.0/24", teamNumber),
					Gateway: formatStr("10.0.%d.254", teamNumber),
				},
			},
		},
	}
}
