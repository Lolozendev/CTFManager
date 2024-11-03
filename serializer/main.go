package serializer

import (
	"fmt"
	"path"
	"strings"

	"github.com/Lolozendev/CTFManager/challengemanager"
	"github.com/Lolozendev/CTFManager/internal"
	"github.com/Lolozendev/CTFManager/internal/constants"
	"github.com/Lolozendev/CTFManager/internal/model"
	"github.com/Lolozendev/CTFManager/internal/model/network"
	"github.com/Lolozendev/CTFManager/internal/model/services"
	"gopkg.in/yaml.v3"
)

var logger = internal.GetLogger()

func SerializeTeam(teamNumber int, teamName string, teamMembers []string) (string, error) {
	teamStrucure, err := makeTeamStructure(teamNumber, teamName, teamMembers)
	if err != nil {
		logger.Errorf("Error: Cannot generate team structure")
		return "", err
	}
	yamlSection, err := makeTeamSection(teamStrucure)
	if err != nil {
		logger.Errorf("Error: Cannot generate Yaml template")
		return "", err
	}
	challengesStructure, err := makeChallengesStructure(teamName, teamNumber)
	if err != nil {
		logger.Errorf("Error: Cannot generate challenges structure")
		return "", err
	}
	ChallengesSection, err := makeChallengeSection(challengesStructure, teamName)
	if err != nil {
		logger.Errorf("Error: Cannot generate challenges section")
		return "", err
	}
	//merge the two sections
	teamYaml := yamlSection + "\n" + ChallengesSection

	return teamYaml, nil

}

func makeChallengeSection(challenges []services.Challenge, teamName string) (string, error) {
	var challengesSection strings.Builder

	for _, challenge := range challenges {
		challengeBytes, err := yaml.Marshal(&challenge)
		if err != nil {
			logger.Errorf("Error: %v", err)
			return "", err
		}
		challengeYaml := string(challengeBytes)

		_, err = fmt.Fprintln(&challengesSection, strings.ReplaceAll(challengeYaml, "<TEAM_NAME>", teamName))
		if err != nil {
			logger.Error("Error: Cannot append challenge to challenges section")
			return "", err
		}

		_, err = fmt.Fprintln(&challengesSection, strings.ReplaceAll(challengeYaml, "<CHALLENGE_NAME>", challenge.Name))
		if err != nil {
			logger.Error("Error: Cannot append challenge to challenges section")
			return "", err
		}
	}

	challengesYaml := indentText(challengesSection.String(), 4)

	return challengesYaml, nil
}

func makeTeamSection(team *model.Team) (string, error) {
	teamBytes, err := yaml.Marshal(&team)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return "", err
	}
	teamYaml := string(teamBytes)

	teamYaml = strings.ReplaceAll(teamYaml, "<TEAM_NAME>", teamYaml)

	return teamYaml, nil
}

func makeChallengesStructure(teamName string, teamNumber int) ([]services.Challenge, error) {
	// if !challengemanager.CheckChallengeDirectory() {
	// 	logger.Errorf("Error: Challenge directory is malformed")
	// 	return nil, errors.New("challenge directory is malformed")
	// }

	challenges := []services.Challenge{}

	for _, challenge := range challengemanager.GetActivatedChallenges() {
		challengeNumber, challengeName, _ := strings.Cut(challenge, "-")

		challenges = append(challenges, services.Challenge{
			Name:          challengeName,
			Build:         path.Join(constants.ChallengesPath, challenge),
			ContainerName: fmt.Sprintf("%s-%s", teamName, challengeName),
			EnvFile:       path.Join(constants.ChallengesPath, challenge, ".env"),
			Networks: services.ChallengeNetworks{
				TeamNetwork: services.ChallengeTeamNetwork{
					Ipv4Address: fmt.Sprintf("10.0.%d.%s", teamNumber, challengeNumber),
				},
			},
		})
	}
	return challenges, nil
}

func makeTeamStructure(teamNumber int, teamName string, teamMembers []string) (*model.Team, error) {
	// if !teammanager.CheckteamDirectory() {
	// 	logger.Errorf("Error: Challenge directory is malformed")
	// 	return nil, errors.New("challenge directory is malformed")
	// }

	members := []model.Member{}
	for _, member := range teamMembers {
		members = append(members, model.Member{
			Username: member,
		})
	}

	team := &model.Team{
		Name:    teamName,
		Number:  teamNumber,
		Members: members,
		Network: network.Network{
			TeamNetwork: network.TeamNetwork{
				Driver: "bridge",
				Ipam: network.Ipam{
					Config: network.Config{
						Subnet: fmt.Sprintf("10.0.%d.0/24", teamNumber),
					},
				},
			},
		},
		Services: services.Services{
			Wireguard: services.Wireguard{
				Image:         "linuxserver/wireguard",
				ContainerName: fmt.Sprintf("%s-wireguard", teamName),
				Ports:         []string{fmt.Sprintf("50%3d:51820/udp", teamNumber)},
			},
			Dnsmasq: services.Dnsmasq{
				Image:         "strm/dnsmasq",
				ContainerName: fmt.Sprintf("%s-dnsmasq", teamName),
				Volumes:       []string{"./dns/dnsmasq.conf:/etc/dnsmasq.conf"},
				Network: services.DnsmasqNetworks{
					TeamNetwork: services.DnsmasqTeamNetwork{
						Ipv4Address: fmt.Sprintf("10.0.%d.253", teamNumber),
					},
				},
			},
		},
	}
	return team, nil

}

func indentText(text string, spaces int) string {
	indent := strings.Repeat(" ", spaces)

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
