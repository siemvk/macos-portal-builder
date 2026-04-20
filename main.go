package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConfigType struct {
	GameToBuild string

	skipCleanup       bool
	showCommandOutput bool
	dryRun            bool
	skipBuild         bool

	tempRepoDir string
	repoUrl     string
}

var Config = ConfigType{
	GameToBuild: "",
	skipCleanup: false,
	dryRun:      false,
	skipBuild:   false,

	tempRepoDir: os.ExpandEnv("$HOME/.temp/source-engine-build-tool"),
	repoUrl:     "https://github.com/nillerusr/source-engine",
}

var ARE_WE_BUILDING_TO_A_APP = true

// logging shit

type loggerType struct {
	debugMsg   func(message string)
	errorMsg   func(message string)
	successMsg func(message string)
	infoMsg    func(message string)
	warnMsg    func(message string)
}

var logLevel = 2

var logger = loggerType{
	debugMsg: func(message string) {
		if logLevel >= 3 {
			fmt.Printf("\033[37m[DEBUG] %s\033[0m\n", message)
		}
	},
	errorMsg: func(message string) {
		if logLevel >= 0 {
			fmt.Printf("\033[31m[ERROR]\033[0m %s\n", message)
		}
	},
	successMsg: func(message string) {
		if logLevel >= 2 {
			fmt.Printf("\033[32m[SUCCESS]\033[0m %s\n", message)
		}
	},
	infoMsg: func(message string) {
		if logLevel >= 2 {
			fmt.Printf("\033[34m[INFO]\033[0m %s\n", message)
		}
	},
	warnMsg: func(message string) {
		if logLevel >= 1 {
			fmt.Printf("\033[33m[WARNING]\033[0m %s\n", message)
		}
	},
}

// returnt of het gelukt is
func execSafe(command string) bool {
	logger.debugMsg("Running command " + command)
	var cmd *exec.Cmd
	if !Config.showCommandOutput {
		cmd = exec.Command("bash", "-c", command+">&/dev/null")
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	err := cmd.Run()
	return err == nil
}

func shellQuote(input string) string {
	return "'" + strings.ReplaceAll(input, "'", "'\\''") + "'"
}

func cleanupTempRepo() {
	if Config.skipCleanup {
		logger.warnMsg("Skipping cleanup!")
		return
	}

	logger.infoMsg("Cleaning up temporary repository directory...")
	if err := os.RemoveAll(Config.tempRepoDir); err != nil {
		logger.errorMsg("Failed to clean up temporary repository directory: " + err.Error())
		return
	}
	logger.successMsg("done cleaning up temporary repository directory!")
}

func gameNameToDir(gameName string) string {
	homeDir := os.ExpandEnv("$HOME")
	if strings.EqualFold(gameName, "hl2") {
		return homeDir + "/Library/Application Support/Steam/steamapps/common/Half-Life 2"
	}
	return homeDir + "/Library/Application Support/Steam/steamapps/common/Portal"
}

func normalizeGameName(gameName string) string {
	return strings.ToLower(strings.TrimSpace(gameName))
}

func validateGameName(gameName string) bool {
	switch normalizeGameName(gameName) {
	case "portal", "hl2":
		return true
	default:
		return false
	}
}

func build() bool {
	logger.debugMsg("starting build process for game: " + Config.GameToBuild)

	logger.infoMsg("Cloning the repo....")
	execSafe("git clone --recursive " + Config.repoUrl + " " + Config.tempRepoDir)
	logger.successMsg("Done cloning repo")

	logger.infoMsg("Installing dependencies...")
	logger.debugMsg("Using Homebrew to install dependencies. This may take a while...")
	execSafe("brew install sdl2 freetype2 fontconfig pkg-config opus jpeg jpeg-turbo libpng libedit")
	logger.debugMsg("Installing Xcode Command Line Tools. This may take a while...")
	execSafe("xcode-select --install")
	logger.successMsg("done installing dependencies!")

	logger.infoMsg("Configuring build script...")

	try1 := execSafe("cd " + Config.tempRepoDir + " && python3 waf configure -T release --prefix='' --build-games=" + Config.GameToBuild)
	if !try1 {
		logger.errorMsg("Basic install failed! This is not uncommon, trying again with different clang")
		try2 := execSafe("cd " + Config.tempRepoDir + " && export CC=/usr/bin/clang && export CXX=/usr/bin/clang++ && python3 waf configure -T release --prefix='' --build-games=" + Config.GameToBuild)
		if !try2 {
			logger.errorMsg("Install failed again! I do not experience this on my machine, so I am doing random fixes from reddit now.")
			try3 := execSafe("cd " + Config.tempRepoDir + " && export CC=/usr/bin/clang && export CXX=/usr/bin/clang++ && arch -arm64 python3 waf configure -T release --prefix='' --build-games=" + Config.GameToBuild)
			if !try3 {
				logger.errorMsg("Install failed again!!!! Okay so what if the first fix broke the second fix so lets try the second fix without the first fix.")
				try4 := execSafe("cd " + Config.tempRepoDir + " && arch -arm64 python3 waf configure -T release --prefix='' --build-games=" + Config.GameToBuild)
				if !try4 {
					logger.errorMsg("Install failed again!!!!! I give up. Please open an issue with the log output so I can try to fix this.")
					cleanupTempRepo()
					logger.errorMsg("Open a issue!!!!")
					return false
				}
			}
		}
	}
	logger.successMsg("done configuring build script!")

	logger.infoMsg("Building the game.... This may take a while...")
	if !Config.skipBuild {
		execSafe("cd " + Config.tempRepoDir + " && python3 waf build")
		logger.successMsg("done building the game!")
	} else {
		logger.warnMsg("Skipping build process!")
	}

	logger.infoMsg("Installing the game to a temp directory...")
	if !execSafe("cd " + Config.tempRepoDir + " && python3 waf install --destdir=" + shellQuote(Config.tempRepoDir+"/installingthismf")) {
		logger.errorMsg("Failed to install build artifacts to temporary directory")
		cleanupTempRepo()
		return false
	}

	logger.successMsg("done installing the game!")

	if Config.dryRun {
		logger.warnMsg("Dry run enabled, skipping installation to game folder.")
		cleanupTempRepo()
		return true
	}

	gameDir := gameNameToDir(Config.GameToBuild)
	if err := os.MkdirAll(gameDir, 0755); err != nil {
		logger.errorMsg("Failed to create game directory: " + err.Error())
		cleanupTempRepo()
		return false
	}

	logger.infoMsg("copying files to the game folder...")
	logger.debugMsg("Copying files from " + Config.tempRepoDir + "/installingthismf to " + gameDir)
	copyCmd := "cd " + shellQuote(gameDir) +
		" && rm -rf ./portal/bin ./bin" +
		" && cp -r " + shellQuote(Config.tempRepoDir+"/installingthismf/portal/bin") + " ./portal/bin" +
		" && cp -r " + shellQuote(Config.tempRepoDir+"/installingthismf/bin") + " ./bin" +
		" && mv ./hl2_osx ./hl2_osx_backup" +
		" && mv " + shellQuote(Config.tempRepoDir+"/installingthismf/hl2_launcher") + " ./hl2_osx"

	if !execSafe(copyCmd) {
		logger.errorMsg("Failed while copying files into game directory")
		cleanupTempRepo()
		return false
	}

	logger.successMsg("done copying files to the game folder!")
	cleanupTempRepo()
	return true
}

func main() {

	repoUrlInput := flag.String("url", "https://github.com/nillerusr/source-engine", "The url of the modified source engine repo.")
	gameBuildInput := flag.String("game", "", "The game to build. Options are: portal and hl2 I can't test hl2 (I don't have it) but it should work, if it doesn't please open an issue.")
	loggerlvlInput := flag.Int("log-level", 2, "0 = only error, 1 = error + warn, 2 = info, success, warn and error, 3 everything")
	testStuff := flag.Bool("testing", false, "Overwrite the config with the one for testing and do some other stuff")
	skipCleanupInput := flag.Bool("skip-cleanup", false, "Whether to skip the cleanup process (deleting the temp repo folder)")
	skipBuildInput := flag.Bool("skip-build", false, "Whether to skip the build process (for testing purposes)")
	tempRepoDirInput := flag.String("temp-repo-dir", os.ExpandEnv("$HOME/.temp/source-engine-build-tool"), "The directory to clone the repo to and build in. By default this is $HOME/.temp/source-engine-build-tool")

	flag.Parse()

	Config.repoUrl = *repoUrlInput
	logLevel = *loggerlvlInput
	Config.skipCleanup = *skipCleanupInput
	Config.skipBuild = *skipBuildInput
	Config.GameToBuild = normalizeGameName(*gameBuildInput)
	Config.tempRepoDir = *tempRepoDirInput
	Config.showCommandOutput = logLevel >= 3

	// If game is not provided via flags and we are in app mode (or just want interactive), ask for it
	if Config.GameToBuild == "" {
		logger.infoMsg("What game do you want to build? (portal/hl2)")
		var userInput string
		fmt.Scanln(&userInput)
		Config.GameToBuild = normalizeGameName(userInput)
		if Config.GameToBuild != "" {
			logger.infoMsg("Good choice! Building " + Config.GameToBuild + " now!")
		}
	}

	if !validateGameName(Config.GameToBuild) {
		logger.errorMsg("Unsupported game: " + *gameBuildInput + ". Supported values are portal and hl2.")
		os.Exit(1)
	}

	if *testStuff {
		logLevel = 3
		logger.debugMsg("debug")
		logger.infoMsg("info")
		logger.warnMsg("warn")
		logger.successMsg("success")
		logger.errorMsg("error D:")

		Config.skipCleanup = true

	}

	success := build()

	if success && !Config.dryRun {
		logger.infoMsg("The game has been successfully built and installed!")
		logger.infoMsg("Would you like to delete this builder tool now? (y/n)")
		var deleteChoice string
		fmt.Scanln(&deleteChoice)
		deleteChoice = strings.ToLower(strings.TrimSpace(deleteChoice))
		if deleteChoice == "y" || deleteChoice == "yes" {
			selfDelete()
		}
	}
}

func selfDelete() {
	executable, err := os.Executable()
	if err != nil {
		logger.errorMsg("Failed to get executable path: " + err.Error())
		return
	}

	// If we are in an app bundle, we might want to delete the whole .app
	// Path is typically .../Contents/MacOS/binary
	// We use filepath.Dir to go up from binary -> MacOS -> Contents -> .app
	if strings.Contains(executable, ".app/Contents/MacOS/") {
		// Better way to find the .app bundle root:
		dir := executable
		for {
			if strings.HasSuffix(dir, ".app") {
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir { // root reached
				break
			}
			dir = parent
		}

		if strings.HasSuffix(dir, ".app") {
			logger.infoMsg("Deleting app bundle: " + dir)
			cmd := exec.Command("rm", "-rf", dir)
			err = cmd.Start()
			if err != nil {
				logger.errorMsg("Failed to start deletion command: " + err.Error())
			}
			return
		}
	}

	// Fallback or non-app binary
	logger.infoMsg("Deleting executable: " + executable)
	// On Unix we can just remove it
	err = os.Remove(executable)
	if err != nil {
		logger.errorMsg("Failed to delete executable: " + err.Error())
		// Fallback to rm just in case
		exec.Command("rm", executable).Start()
	}
	logger.successMsg("Cleanup initiated. Goodbye!")
}