package services

type Challenge struct {
	Name          string
	Build         string            `yaml:"build"`
	ContainerName string            `yaml:"{ TEAM_NAME }-{ CHALLENGE_NAME }"`
	EnvFile       string            `yaml:"env_file,omitempty"`
	Networks      ChallengeNetworks `yaml:"networks"`
}

type ChallengeNetworks struct {
	TeamNetwork ChallengeTeamNetwork `yaml:"{ TEAM_NAME }-Network"`
}

type ChallengeTeamNetwork struct {
	Ipv4Address string `yaml:"ipv4_address"`
}
