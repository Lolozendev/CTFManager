package services

type Dnsmasq struct {
	Image         string          `yaml:"image"`
	ContainerName string          `yaml:"container_name"`
	Volumes       []string        `yaml:"volumes"`
	Network       DnsmasqNetworks `yaml:"networks"`
}

type DnsmasqNetworks struct {
	TeamNetwork DnsmasqTeamNetwork `yaml:"{ TEAM_NAME }-Network"`
}

type DnsmasqTeamNetwork struct {
	Ipv4Address string `yaml:"ipv4_address"`
}

/*
dnsmasq:
	image: strm/dnsmasq
	volumes:
	- ./dns/dnsmasq.conf:/etc/dnsmasq.conf
	cap_add:
	- NET_ADMIN
	networks:
		team_network:
			ipv4_address: 10.0.{ TEAM_NUMBER }.253
*/
