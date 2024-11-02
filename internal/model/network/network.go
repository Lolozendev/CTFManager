package network

type Network struct {
	TeamNetwork TeamNetwork `yaml:"{ TEAM_NAME }-Network"`
}

type TeamNetwork struct {
	Driver string `yaml:"driver"`
	Ipam   Ipam   `yaml:"ipam"`
}

type Ipam struct {
	Config Config `yaml:"config"`
}

type Config struct {
	Subnet string `yaml:"subnet"`
}

/*
networks:
  { TEAM_NAME }-Network:
    driver: bridge
	ipam:
	  config:
*/
