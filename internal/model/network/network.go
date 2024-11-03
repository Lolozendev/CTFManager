package network

type Network struct {
	TeamNetwork TeamNetwork `yaml:"<TEAM_NAME>-Network"`
}

type TeamNetwork struct {
	Driver string `yaml:"driver"`
	Ipam   Ipam   `yaml:"ipam"`
}

type Ipam struct {
	Config []Config `yaml:"config"`
}

type Config struct {
	Subnet  string `yaml:"subnet"`
	Gateway string `yaml:"gateway"`
}

/*
networks:
	<TEAM_NAME>-Network:
    	driver: bridge
		ipam:
	  		config: 10.0.<TEAM_NUMBER>.0/24
*/
