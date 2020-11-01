package core

import (
	"fmt"
	"runtime"
)

const (
	AppName          string = "p2pNG-core"
	AppIntroduction  string = "Next Generation Peer-to-Peer Network Core."
	LongIntroduction string = "p2pNG-core is universal platform for peer to peer network.\n" +
		"It is going to implements all the basic components of our p2pNG Protocol."
)

var (
	buildVersionTag = "0.0.0-dev"
	buildTime       = ""
	buildName       = "Common Dev"
)

func GetAppStatement() string {
	sTime := ""
	//noinspection GoBoolExpressions
	if GetBuildTime() != "" {
		sTime = " " + GetBuildTime() + ","
	}
	return fmt.Sprintf("%s, %s v%s (%s Build,%s %s/%s)", AppName, AppIntroduction, GetVersionTag(),
		buildName, sTime, runtime.GOOS, runtime.GOARCH)
}
func GetVersionStatement() string {
	sTime := ""
	//noinspection GoBoolExpressions
	if GetBuildTime() != "" {
		sTime = " " + GetBuildTime() + ","
	}
	return fmt.Sprintf("v%s (%s Build,%s %s/%s)", GetVersionTag(),
		buildName, sTime, runtime.GOOS, runtime.GOARCH)
}
func GetVersionTag() string { return buildVersionTag }
func GetBuildTime() string  { return buildTime }
func GetBuildName() string  { return buildName }
