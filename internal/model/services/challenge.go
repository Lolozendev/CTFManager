package services

type Challenge struct {
	Name          string            `yaml:"-"`
	Build         string            `yaml:"build"`
	ContainerName string            `yaml:"container_name"`
	EnvFile       string            `yaml:"env_file,omitempty"`
	Networks      ChallengeNetworks `yaml:"networks"`
}

type ChallengeNetworks struct {
	TeamNetwork ChallengeTeamNetwork `yaml:"<TEAM_NAME>-Network"`
}

type ChallengeTeamNetwork struct {
	Ipv4Address string `yaml:"ipv4_address"`
}

/*
<CHALLENGE_NAME>:
	build: /challenges/<CHALLENGE_NUMBER>-<CHALLENGE_NAME>
	container_name: <TEAM_NAME>-<CHALLENGE_NAME>
	env_file: /challenges/<CHALLENGE_NUMBER>-<CHALLENGE_NAME>/.env
	networks:
		TEAM_NAME-Network:
			ipv4_address: 10.0.<TEAM_NUMBER>.<CHALLENGE_NUMBER>
*/
