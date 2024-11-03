// This doesn't returns any value, it just logs the errors and warnings because in a near future , it will handle the errors and warnings on its own.
package teammanager

import (
	"errors"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Lolozendev/CTFManager/internal"
	"github.com/Lolozendev/CTFManager/internal/constants"
	"go.uber.org/zap"
)

const (
	TeamsPath                 string = constants.TeamsPath
	TeamsFirstNetworkPosition int    = constants.TeamsFirstNetworkPosition
	TeamsLastNetworkPosition  int    = constants.TeamsLastNetworkPosition
)

var (
	TeamsNameRegexp *regexp.Regexp     = regexp.MustCompile(`(?P<DIGIT>\d{1,3}|x)-\w*$`)
	logger          *zap.SugaredLogger = internal.GetLogger()
)

func normalizeteamName(name string) bool {
	foundError := false
	if !TeamsNameRegexp.MatchString(name) {
		logger.Error("Error: ", name, " does not match the expected format")
		if regexp.MustCompile(`\w*`).MatchString(name) {
			logger.Info("Renaming ", name, " to a disabled team")
			newname := "x-" + name
			os.Rename(TeamsPath+"/"+name, TeamsPath+"/"+newname)
		} else {
			logger.Error("Error: cannot rename ", name, " to a disabled team, make sure the team name is alphanumeric only")
			//TODO: Add a way to Normalize the name to remove special characters
			foundError = true
		}
	}
	return foundError
}

func checkteamNames() bool {
	foundError := false
	entries, err := os.ReadDir(TeamsPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	for _, entry := range entries {
		foundError = normalizeteamName(entry.Name())
	}
	return foundError
}

func checkDuplicatesAndHoles() (map[int]string, error) {
	entries, err := os.ReadDir(TeamsPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	actualTeams := make(map[int]string, 0)

	for _, entry := range entries {
		digit, name, _ := strings.Cut(entry.Name(), "-")
		if digit == "x" {
			continue
		}
		digitInt, _ := strconv.Atoi(digit)

		if digitInt < TeamsFirstNetworkPosition || digitInt > TeamsLastNetworkPosition {
			logger.Error("Error: ", digitInt, " is not in the expected range")
			return nil, errors.New("error: team id is not in the expected range")
		}

		//check if there is duplicate Teams ids
		if _, ok := actualTeams[digitInt]; !ok {
			actualTeams[digitInt] = name
		} else {
			logger.Error("Error: Duplicate team name ", name, " found")
			return nil, errors.New("error: Duplicate team names found")
		}
	}

	//check if there are holes in the Teams ids
	for i := TeamsFirstNetworkPosition; i <= TeamsLastNetworkPosition; i++ {
		if _, ok := actualTeams[i]; !ok {
			logger.Error("Error: Id ", i, " is not used")
			return nil, errors.New("error: There are holes in the Teams ids")
		}
	}

	return actualTeams, nil

}

func checkTeamstructure(path string) bool {
	//check if the team directory has a dockerfile and a .env file
	if _, err := os.Stat(path + "/compose.yml"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Error("Error: ", path, " does not have a compose.yml file")
			return false
		} else {
			logger.Error("Unkown error: ", err)
			return false
		}
	}
	return true
}

func checkTeamsStructure() bool {
	foundError := false
	entries, err := os.ReadDir(TeamsPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	for _, entry := range entries {
		if checkTeamstructure(TeamsPath + "/" + entry.Name()) {
			logger.Info("team ", entry.Name(), " has the correct structure")
		} else {
			foundError = true
			logger.Error("Error: team ", entry.Name(), " does not have the correct structure")
		}
	}
	return foundError
}

func noTeams() bool {
	entries, err := os.ReadDir(TeamsPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	if len(entries) == 0 {
		logger.Error("Error: ", TeamsPath, " is empty")
		return true
	}
	return false
}

func CheckteamDirectory() bool {
	fileInfo, err := os.Stat(TeamsPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Error("Error: ", TeamsPath, " does not exist")
		} else {
			logger.Error("Unkown error: ", err)
		}
		return false
	}

	if !fileInfo.IsDir() {
		logger.Error("Error: ", TeamsPath, `is not a directory,
		please create a directory named Teams`)
		return false
	}

	if noTeams() {
		return false
	}

	if checkteamNames() {
		return false
	}

	_, err = checkDuplicatesAndHoles()
	if err != nil {
		return false
	}

	return checkTeamsStructure()
}

func GetTeams() []string {
	entries, err := os.ReadDir(TeamsPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	Teams := make([]string, 0)

	for _, entry := range entries {
		Teams = append(Teams, entry.Name())
	}
	return Teams
}

func GetActivatedTeams() []string {
	Teams := GetTeams()
	activatedTeams := make([]string, 0)
	for _, team := range Teams {
		if !strings.HasPrefix(team, "x-") {
			activatedTeams = append(activatedTeams, team)
		}
	}
	return activatedTeams
}

func GetDisabledTeams() []string {
	Teams := GetTeams()
	disabledTeams := make([]string, 0)
	for _, team := range Teams {
		if strings.HasPrefix(team, "x-") {
			disabledTeams = append(disabledTeams, team)
		}
	}
	return disabledTeams
}
