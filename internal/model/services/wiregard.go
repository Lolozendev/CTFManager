package services

type Wireguard struct {
	ServerAddress string          `yaml:"-"`
	Image         string          `yaml:"image"`
	ContainerName string          `yaml:"container_name"`
	Ports         []string        `yaml:"ports"`
	Environment   []string        `yaml:"environment"`
	Volumes       []string        `yaml:"volumes"`
	CapAdd        []string        `yaml:"cap_add"`
	Network       WiregardNetwork `yaml:"networks"`
}

type WiregardNetwork struct {
	TeamNetwork WiregardTeamNetwork `yaml:"<TEAM_NAME>-Network"`
}

type WiregardTeamNetwork struct {
	Ipv4Address string `yaml:"ipv4_address"`
}

/*
wireguard:
    	image: linuxserver/wireguard
    	container_name: wireguard
    	ports:
      	- "50<TEAM_NUMBER>:51820/udp"
    	environment:
      	- PUID=1000
      	- PGID=1000
      	- TZ=Europe/Paris
      	- PEERS=<TEAM_MEMBERS>
      	- PEERDNS=10.0.<TEAM_NUMBER>.253,1.1.1.1
      	- ALLOWEDIPS=10.0.<TEAM_NUMBER>.0/24
      	- SERVERURL=127.0.0.1 //replace with the server's IP
      	- SERVERPORT=50<TEAM_NUMBER>
    	volumes:
      	- ./config:/config
    	cap_add:
      	- NET_ADMIN
    	networks:
      		team_network:
        		ipv4_address: 10.0.<TEAM_NUMBER>.252
*/
