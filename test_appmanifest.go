package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	content := `"AppState"
{
	"appid"		"400"
	"Universe"		"1"
	"name"		"Portal"
	"StateFlags"		"4"
	"installdir"		"Portal"
	"LastUpdated"		"1711202812"
	"UpdateResult"		"0"
	"SizeOnDisk"		"5234515569"
	"buildid"		"9537153"
	"LastOwner"		"123"
	"BytesToDownload"		"0"
	"BytesDownloaded"		"0"
	"BytesToStage"		"0"
	"BytesStaged"		"0"
	"TargetBuildID"		"9537153"
	"AutoUpdateBehavior"		"0"
	"AllowOtherDownloadsWhileRunning"		"0"
	"ScheduledAutoUpdate"		"0"
	"InstalledDepots"
	{
		"401"
		{
			"manifest"		"1433018244247590897"
			"size"		"5130768997"
		}
	}
	"UserConfig"
	{
		"language"		"english"
	}
}`
	fmt.Println(strings.Contains(content, "betakey"))
}
