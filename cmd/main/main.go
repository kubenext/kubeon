package main

import (
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

}
