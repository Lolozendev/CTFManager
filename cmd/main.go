package main

import (
	"fmt"
	"os"

	"github.com/Lolozendev/CTFManager/internal"
	"github.com/Lolozendev/CTFManager/serializer"
)

func main() {
	logger := internal.GetLogger()
	defer logger.Sync()

	logger.Info("only usable command for now is 'team create <teamname>'")

	if len(os.Args) < 4 {
		logger.Error("Error: Missing arguments")
		return
	}

	logger.Info("creating team ", os.Args[3])

	serialized, err := serializer.SerializeTeam(1, os.Args[3], []string{"user1", "user2"})
	if err != nil {
		logger.Error("Error: ", err)
		return
	}

	fmt.Println(serialized)
}
