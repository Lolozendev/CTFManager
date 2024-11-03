package constants

// This holds the constants for the paths of the challenges and teams used in multiple part of the project
const (
	ChallengesPath                 string = "/challenges"
	ChallengesFirstNetworkPosition int    = 11
	ChallengesLastNetworkPosition  int    = 249
	TeamsPath                      string = "/equipes"
	TeamsFirstNetworkPosition      int    = 1
	TeamsLastNetworkPosition       int    = 254
	DnsMasqTemplatePath            string = "/dnsconf/dnsmasq.template"

	//placeholders
	TeamNamePlaceholder        string = "<TEAM_NAME>"
	TeamNumberPlaceholder      string = "<TEAM_NUMBER>"
	ChallengeNamePlaceholder   string = "<CHALLENGE_NAME>"
	ChallengeNumberPlaceholder string = "<CHALLENGE_NUMBER>"
)
