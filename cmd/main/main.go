package main

import (
	"github.com/kubenext/kubeon/internal/command"
	"log"
	"math/rand"
	"time"
)

var (
	version   = "devVersion"
	gitCommit = "devGitCommit"
	buildTime = "devBuildTime"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	// remove timestamp from log
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func main() {
	command.A()
}
