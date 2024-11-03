package services

type Dnsmasq struct {
	Image         string          `yaml:"image"`
	ContainerName string          `yaml:"container_name"`
	Volumes       []string        `yaml:"volumes"`
	Network       DnsmasqNetworks `yaml:"networks"`
}

type DnsmasqNetworks struct {
	TeamNetwork DnsmasqNetwork `yaml:"<TEAM_NAME>-Network"`
}

type DnsmasqNetwork struct {
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
		<TEAM_NAME>-Network:
			ipv4_address: 10.0.<TEAM_NUMBER>.253
*/
