package main

import (
	"github.com/kubenext/kubethan/internal/commands"
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
}

func main() {
	commands.Execute(version, gitCommit, buildTime)
}
