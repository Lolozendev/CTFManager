package challengemanager

import (
	"errors"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Lolozendev/CTFManager/internal"
	"go.uber.org/zap"
)

const (
	ChallengesPath                 string = "/challenges"
	ChallengesFirstNetworkPosition int    = 11
	ChallengesLastNetworkPosition  int    = 249
)

var (
	ChallengesNameRegexp *regexp.Regexp     = regexp.MustCompile(`(?P<DIGIT>\d{1,3}|x)-\w*$`)
	logger               *zap.SugaredLogger = internal.GetLogger()
)

func normalizeChallengeName(name string) {
	if !ChallengesNameRegexp.MatchString(name) {
		logger.Error("Error: ", name, " does not match the expected format")
		if regexp.MustCompile(`\w*`).MatchString(name) {
			logger.Info("Renaming ", name, " to a disabled challenge")
			var newname = "x-" + name
			os.Rename(ChallengesPath+"/"+name, ChallengesPath+"/"+newname)
		} else {
			logger.Error("Error: cannot rename ", name, " to a disabled challenge, make sure the challenge name is alphanumeric only")
			//TODO: Add a way to Normalize the name to remove special characters
		}
	}
}

func CheckChallengeNames() {
	entries, err := os.ReadDir(ChallengesPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	for _, entry := range entries {
		normalizeChallengeName(entry.Name())
	}
}

func checkDuplicatesAndHoles() (map[int]string, error) {
	entries, err := os.ReadDir(ChallengesPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	var actualChallenges = make(map[int]string, 0)

	for _, entry := range entries {
		digit, name, _ := strings.Cut(entry.Name(), "-")
		if digit == "x" {
			continue
		}
		digitInt, _ := strconv.Atoi(digit)

		if digitInt < ChallengesFirstNetworkPosition || digitInt > ChallengesLastNetworkPosition {
			logger.Error("Error: ", digitInt, " is not in the expected range")
			return nil, errors.New("error: Challenge id is not in the expected range")
		}

		//check if there is duplicate challenges ids
		if _, ok := actualChallenges[digitInt]; !ok {
			actualChallenges[digitInt] = name
		} else {
			logger.Error("Error: Duplicate challenge name ", name, " found")
			return nil, errors.New("error: Duplicate challenge names found")
		}
	}

	//check if there are holes in the challenges ids
	for i := ChallengesFirstNetworkPosition; i <= ChallengesLastNetworkPosition; i++ {
		if _, ok := actualChallenges[i]; !ok {
			logger.Error("Error: Id ", i, " is not used")
			return nil, errors.New("error: There are holes in the challenges ids")
		}
	}

	return actualChallenges, nil

}

func checkChallengeStructure(path string) bool {
	//check if the challenge directory has a dockerfile and a .env file
	if _, err := os.Stat(path + "/Dockerfile"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Error("Error: ", path, " does not have a Dockerfile")
			return false
		} else {
			logger.Error("Unkown error: ", err)
			return false
		}
	}
	return true
}

func CheckChallengesStructure() bool {
	var foundError bool = false
	entries, err := os.ReadDir(ChallengesPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	for _, entry := range entries {
		if checkChallengeStructure(ChallengesPath + "/" + entry.Name()) {
			logger.Info("Challenge ", entry.Name(), " has the correct structure")
		} else {
			foundError = true
			logger.Error("Error: Challenge ", entry.Name(), " does not have the correct structure")
		}
	}
	return foundError
}

func noChallenges() bool {
	entries, err := os.ReadDir(ChallengesPath)
	if err != nil {
		logger.Error("Error: ", err)
	}

	if len(entries) == 0 {
		logger.Error("Error: ", ChallengesPath, " is empty")
		return true
	}
	return false
}

func CheckChallengeDirectory() bool {
	fileInfo, err := os.Stat(ChallengesPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Error("Error: ", ChallengesPath, " does not exist")
		} else {
			logger.Error("Unkown error: ", err)
		}
		return false
	}

	if !fileInfo.IsDir() {
		logger.Error("Error: ", ChallengesPath, `is not a directory,
		please create a directory named challenges`)
		return false
	}

	if noChallenges() {
		return false
	}

	_, err = checkDuplicatesAndHoles()
	if err != nil {
		return false
	}

	return CheckChallengesStructure()
}
