package model

import "fmt"

// Service represents a Docker Compose service configuration
type Service struct {
	Image         string            `yaml:"image,omitempty"`
	Build         string            `yaml:"build,omitempty"`
	ContainerName string            `yaml:"container_name"`
	Ports         []string          `yaml:"ports,omitempty"`
	Environment   []string          `yaml:"environment,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	CapAdd        []string          `yaml:"cap_add,omitempty"`
	EnvFile       string            `yaml:"env_file,omitempty"`
	Networks      map[string]IPAddr `yaml:"networks"`
}

// IPAddr represents network IP configuration
type IPAddr struct {
	Ipv4Address string `yaml:"ipv4_address"`
}

// NewWireguardService creates a Wireguard VPN service
func NewWireguardService(teamName string, teamNumber int, memberCount int) Service {
	networkName := teamName + "-Network"
	return Service{
		Image:         "linuxserver/wireguard",
		ContainerName: teamName + "-wireguard",
		Ports:         []string{formatPort(50000 + teamNumber)},
		Environment: []string{
			"PUID=1000",
			"PGID=1000",
			"TZ=Europe/Paris",
			formatEnv("PEERS", memberCount),
			formatEnv("PEERDNS", formatIP(teamNumber, 253)),
			formatEnv("ALLOWEDIPS", formatSubnet(teamNumber)),
			"SERVERURL=127.0.0.1", // TODO: Make configurable
			formatEnv("SERVERPORT", 50000+teamNumber),
		},
		Volumes: []string{"./config:/config"},
		CapAdd:  []string{"NET_ADMIN"},
		Networks: map[string]IPAddr{
			networkName: {Ipv4Address: formatIP(teamNumber, 252)},
		},
	}
}

// NewDnsmasqService creates a DNS service
func NewDnsmasqService(teamName string, teamNumber int) Service {
	networkName := teamName + "-Network"
	return Service{
		Image:         "strm/dnsmasq",
		ContainerName: teamName + "-dnsmasq",
		Volumes:       []string{"./dns/dnsmasq.conf:/etc/dnsmasq.conf"},
		Networks: map[string]IPAddr{
			networkName: {Ipv4Address: formatIP(teamNumber, 253)},
		},
	}
}

// NewChallengeService creates a challenge service
func NewChallengeService(teamName string, teamNumber int, challengeNumber int, challengeName string, buildPath string, envPath string) Service {
	networkName := teamName + "-Network"
	return Service{
		Build:         buildPath,
		ContainerName: teamName + "-" + challengeName,
		EnvFile:       envPath,
		Networks: map[string]IPAddr{
			networkName: {Ipv4Address: formatIP(teamNumber, challengeNumber)},
		},
	}
}

// Helper functions
func formatPort(port int) string {
	return formatStr("%d:51820/udp", port)
}

func formatEnv(key string, value interface{}) string {
	return formatStr("%s=%v", key, value)
}

func formatIP(teamNumber, host int) string {
	return formatStr("10.0.%d.%d", teamNumber, host)
}

func formatSubnet(teamNumber int) string {
	return formatStr("10.0.%d.0/24", teamNumber)
}

func formatStr(format string, args ...interface{}) string {
	// Simple helper to avoid importing fmt in every function
	return fmt.Sprintf(format, args...)
}
