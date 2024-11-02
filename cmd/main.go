package main

import (
	"github.com/Lolozendev/CTFManager/internal"
)

func main() {
	logger := internal.GetLogger()
	defer logger.Sync()
	logger.Info("Hello, World!")
	logger.Warn("This is a warning!")
	logger.Error("This is an error!")
}
