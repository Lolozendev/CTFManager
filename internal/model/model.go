package model

type Team struct {
	Networks Networks `yaml:"networks"`
	Services Services `yaml:"services"`
	Name     string
}

type Member struct {
	Username string
}

type Networks struct {
	TeamNetwork struct {
		Driver string `yaml:"driver"`
		Ipam   struct {
			Config []struct {
				Subnet string `yaml:"subnet"`
			} `yaml:"config"`
		} `yaml:"ipam"`
	} `yaml:"{ TEAM_NAME }-Network"`
}

type Services struct {
	Wireguard Wireguard `yaml:"wireguard"`
	Dnsmasq   Dnsmasq   `yaml:"dnsmasq"`
	Challenge []Challenge
}

type Wireguard struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Ports         []string `yaml:"ports"`
	Environment   []string `yaml:"environment"`
	Volumes       []string `yaml:"volumes"`
	CapAdd        []string `yaml:"cap_add"`
	Networks      struct {
		TeamNetwork struct {
			Ipv4Address string `yaml:"ipv4_address"`
		} `yaml:"{ TEAM_NAME }-Network"`
	} `yaml:"networks"`
}

type Dnsmasq struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Volumes       []string `yaml:"volumes"`
	Networks      struct {
		TeamNetwork struct {
			Ipv4Address string `yaml:"ipv4_address"`
		} `yaml:"{ TEAM_NAME }-Network"`
	} `yaml:"networks"`
}

type Challenge struct {
	Name          string
	Build         string `yaml:"build"`
	ContainerName string `yaml:"{ TEAM_NAME }-{ CHALLENGE_NAME }"`
	EnvFile       string `yaml:"env_file,omitempty"`
	Networks      struct {
		TeamNetwork struct {
			Ipv4Address string `yaml:"ipv4_address"`
		} `yaml:"{ TEAM_NAME }-Network"`
	} `yaml:"networks"`
}
