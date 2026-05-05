package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
func execSafe(args ...string) bool {
	logger.debugMsg("Running command " + strings.Join(args, " "))
	if len(args) == 0 {
		return false
	}
	cmd := exec.Command(args[0], args[1:]...)
	if Config.showCommandOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	return err == nil
}

func execSafeDirEnv(dir string, env []string, args ...string) bool {
	logger.debugMsg("Running command in " + dir + " with env " + strings.Join(env, " ") + ": " + strings.Join(args, " "))
	if len(args) == 0 {
		return false
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	if Config.showCommandOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	return err == nil
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
	logger.successMsg("Done cleaning up temporary repository directory!")
}

var (
	cachedSteamLibraries []string
	steamLibrariesOnce   sync.Once
)

func findSteamLibraries() []string {
	steamLibrariesOnce.Do(func() {
		homeDir := os.ExpandEnv("$HOME")
		defaultSteamPath := filepath.Join(homeDir, "Library", "Application Support", "Steam")
		libraries := []string{defaultSteamPath}

		vdfPath := filepath.Join(defaultSteamPath, "steamapps", "libraryfolders.vdf")
		content, err := os.ReadFile(vdfPath)
		if err == nil {
			pathRegex := regexp.MustCompile(`(?i)"path"\s+"([^"]+)"`)
			matches := pathRegex.FindAllStringSubmatch(string(content), -1)

			seen := make(map[string]bool, len(matches)+len(libraries))
			for _, l := range libraries {
				seen[l] = true
			}

			for _, match := range matches {
				if len(match) == 2 {
					path := match[1]
					path = filepath.Clean(path)
					if !seen[path] {
						seen[path] = true
						libraries = append(libraries, path)
					}
				}
			}
		}
		cachedSteamLibraries = libraries
	})

	result := make([]string, len(cachedSteamLibraries))
	copy(result, cachedSteamLibraries)
	return result
}

func getGameLibraryPath(appId string) string {
	libraries := findSteamLibraries()
	for _, lib := range libraries {
		manifestPath := filepath.Join(lib, "steamapps", fmt.Sprintf("appmanifest_%s.acf", appId))
		if _, err := os.Stat(manifestPath); err == nil {
			return lib // found the library where this game is installed
		}
	}
	return ""
}

func gameNameToDir(gameName string) string {
	var appId string
	var folderName string
	if normalizeGameName(gameName) == "portal" {
		appId = "400"
		folderName = "Portal"
	} else if normalizeGameName(gameName) == "hl2" {
		appId = "220"
		folderName = "Half-Life 2"
	} else {
		return ""
	}

	libPath := getGameLibraryPath(appId)
	if libPath != "" {
		return filepath.Join(libPath, "steamapps", "common", folderName)
	}

	// fallback to default
	homeDir := os.ExpandEnv("$HOME")
	return filepath.Join(homeDir, "Library", "Application Support", "Steam", "steamapps", "common", folderName)
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

func checkSteamBetaRequirement(gameName string) bool {
	var appId string
	if normalizeGameName(gameName) == "portal" {
		appId = "400"
	} else if normalizeGameName(gameName) == "hl2" {
		appId = "220"
	} else {
		return true
	}

	libPath := getGameLibraryPath(appId)
	if libPath == "" {
		logger.errorMsg("Could not find Steam appmanifest for the game.")
		logger.errorMsg("Make sure the game is installed via Steam and located in a valid Steam library. Pirated copies of the game are not supported!")
		return false
	}

	manifestPath := filepath.Join(libPath, "steamapps", fmt.Sprintf("appmanifest_%s.acf", appId))
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		logger.errorMsg("Could not read Steam appmanifest for the game.")
		return false
	}

	contentStr := strings.ToLower(string(content))
	hasBetaKey := strings.Contains(contentStr, "betakey")
	isPublic, _ := regexp.MatchString(`"betakey"\s+"public"`, contentStr)

	isValid := false

	switch appId {
	case "220":
		isSteamLegacy, _ := regexp.MatchString(`"betakey"\s+"steam_legacy"`, contentStr)
		isValid = isSteamLegacy
	case "400":
		isBeta, _ := regexp.MatchString(`"betakey"\s+"beta"`, contentStr)
		isValid = isBeta
	default:
		isValid = (hasBetaKey && !isPublic)
	}

	if !isValid {
		logger.errorMsg("You are on the wrong branch of the game!")
		switch appId {
		case "220":
			logger.errorMsg("For Half-Life 2, you MUST install the 'steam_legacy - Pre-20th Anniversary Build'!")
			logger.errorMsg("Please go to Steam -> right click the game -> Properties -> Betas -> select 'steam_legacy'.")
		case "400":
			logger.errorMsg("For Portal, you MUST install the 'beta - SteamPipe Beta'!")
			logger.errorMsg("Please go to Steam -> right click the game -> Properties -> Betas -> select 'beta'.")
		default:
			logger.errorMsg("You have to have the steam beta or legacy version of the game!")
			logger.errorMsg("Please go to Steam -> right click the game -> Properties -> Betas -> select the beta or legacy branch.")
		}
		return false
	}

	logger.successMsg("Beta/Legacy check passed")
	return true
}

func build() bool {
	logger.debugMsg("Starting build process for game: " + Config.GameToBuild)

	logger.infoMsg("Checking system requirements...")
	xcodeOut, err := exec.Command("xcode-select", "-p").Output()
	hasXcode := false
	if err == nil {
		xcodePath := strings.TrimSpace(string(xcodeOut))
		clangPath := filepath.Join(xcodePath, "usr", "bin", "clang")
		if _, statErr := os.Stat(clangPath); statErr == nil {
			hasXcode = true
		}
	}

	if !hasXcode {
		logger.errorMsg("Xcode Command Line Tools are missing or broken!")
		logger.errorMsg("They are required to compile the game and provide 'git'.")
		logger.errorMsg("Please run this command in your terminal:")
		logger.errorMsg("  xcode-select --install")
		logger.errorMsg("A window will pop up. Follow the prompts to install, wait for it to finish completely, and then run this builder tool again.")
		return false
	}

	brewPath, err := exec.LookPath("brew")
	if err != nil {
		// fallback for if users just installed it but haven't restarted their terminal
		// and intel path and apple silicon path
		if _, err2 := os.Stat("/opt/homebrew/bin/brew"); err2 == nil {
			brewPath = "/opt/homebrew/bin/brew"
		} else if _, err3 := os.Stat("/usr/local/bin/brew"); err3 == nil {
			brewPath = "/usr/local/bin/brew"
		}
	}

	if brewPath == "" {
		logger.errorMsg("Homebrew is not installed! It is required to install build dependencies.")
		logger.errorMsg("Please install Homebrew by running the following command in your terminal:")
		logger.errorMsg(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
		logger.errorMsg("After installing Homebrew, please restart your terminal and run this tool again.")
		return false
	} else if err != nil {
		// brew found via fallbacks o inject it into PATH for the session (why did they not restart the terminal grrr)
		os.Setenv("PATH", filepath.Dir(brewPath)+":"+os.Getenv("PATH"))
	}

	logger.infoMsg("Installing dependencies...")
	logger.debugMsg("Using Homebrew to install dependencies. This may take a while...")
	if !execSafe("brew", "install", "python", "sdl2", "python3", "freetype2", "fontconfig", "pkg-config", "opus", "jpeg", "jpeg-turbo", "libpng", "libedit") {
		logger.warnMsg("Dependencies installation warning. If the build fails later, this might be why.")
	}
	logger.successMsg("Done installing dependencies!")

	logger.infoMsg("Cloning the repo....")
	if !execSafe("git", "clone", "--recursive", Config.repoUrl, Config.tempRepoDir) {
		logger.errorMsg("Failed to clone the repository! Please check your internet connection or git permissions.")
		return false
	}
	logger.successMsg("Done cloning repo")

	logger.infoMsg("Configuring build script...")

	try1 := execSafeDirEnv(Config.tempRepoDir, []string{"CXXFLAGS=-include alloca.h"}, "python3", "waf", "configure", "-T", "release", "--prefix=", "--build-games="+Config.GameToBuild)
	if !try1 {
		logger.errorMsg("Basic install failed! This is not uncommon, trying again with different clang")
		try2 := execSafeDirEnv(Config.tempRepoDir, []string{"CC=/usr/bin/clang", "CXX=/usr/bin/clang++", "CXXFLAGS=-include alloca.h"}, "python3", "waf", "configure", "-T", "release", "--prefix=", "--build-games="+Config.GameToBuild)
		if !try2 {
			logger.errorMsg("Install failed again! I do not experience this on my machine, so I am doing random fixes from reddit now.")
			try3 := execSafeDirEnv(Config.tempRepoDir, []string{"CC=/usr/bin/clang", "CXX=/usr/bin/clang++", "CXXFLAGS=-include alloca.h"}, "arch", "-arm64", "python3", "waf", "configure", "-T", "release", "--prefix=", "--build-games="+Config.GameToBuild)
			if !try3 {
				logger.errorMsg("Install failed again!!!! Okay so what if the first fix broke the second fix so lets try the second fix without the first fix.")
				try4 := execSafeDirEnv(Config.tempRepoDir, []string{"CXXFLAGS=-include alloca.h"}, "arch", "-arm64", "python3", "waf", "configure", "-T", "release", "--prefix=", "--build-games="+Config.GameToBuild)
				if !try4 {
					logger.errorMsg("Install failed again!!!!! I give up. Please open an issue with the log output and device specs so I can try to fix this.")
					cleanupTempRepo()
					logger.errorMsg("Open a issue!!!!")
					return false
				}
			}
		}
	}
	logger.successMsg("Done configuring build script!")

	logger.infoMsg("Building the game.... this may take a while...")
	if !Config.skipBuild {
		if !execSafeDirEnv(Config.tempRepoDir, nil, "python3", "waf", "build") {
			logger.errorMsg("Failed to build the game! Please run with --log-level 3 to see the compile errors.")
			cleanupTempRepo()
			return false
		}
		logger.successMsg("Done building the game!")
	} else {
		logger.warnMsg("Skipping build process!")
	}

	logger.infoMsg("Installing the game to a temp directory...")
	if !execSafeDirEnv(Config.tempRepoDir, nil, "python3", "waf", "install", "--destdir="+Config.tempRepoDir+"/installingthismf") {
		logger.errorMsg("Failed to install build artifacts to temporary directory")
		cleanupTempRepo()
		return false
	}

	logger.successMsg("Done installing the game!")

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

	logger.infoMsg("Copying files to the game folder...")
	logger.debugMsg("Copying files from " + Config.tempRepoDir + "/installingthismf to " + gameDir)

	if !execSafeDirEnv(gameDir, nil, "rm", "-rf", "./"+Config.GameToBuild+"/bin", "./bin") {
		logger.errorMsg("Failed while cleaning game directory")
		cleanupTempRepo()
		return false
	}
	if !execSafeDirEnv(gameDir, nil, "cp", "-r", Config.tempRepoDir+"/installingthismf/"+Config.GameToBuild+"/bin", "./"+Config.GameToBuild+"/bin") {
		logger.errorMsg("Failed while copying game bin into game directory")
		cleanupTempRepo()
		return false
	}
	if !execSafeDirEnv(gameDir, nil, "cp", "-r", Config.tempRepoDir+"/installingthismf/bin", "./bin") {
		logger.errorMsg("Failed while copying bin into game directory")
		cleanupTempRepo()
		return false
	}
	execSafeDirEnv(gameDir, nil, "mv", "./hl2_osx", "./hl2_osx_backup") // Ignore failure
	if !execSafeDirEnv(gameDir, nil, "mv", Config.tempRepoDir+"/installingthismf/hl2_launcher", "./hl2_osx") {
		logger.errorMsg("Failed while copying hl2_launcher into game directory")
		cleanupTempRepo()
		return false
	}

	logger.successMsg("Done copying files to the game folder!")
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
	Config.showCommandOutput = logLevel >= 3

	if err := os.MkdirAll(*tempRepoDirInput, 0755); err != nil {
		logger.errorMsg("Failed to create temporary repository base directory: " + err.Error())
		os.Exit(1)
	}

	tempDir, err := os.MkdirTemp(*tempRepoDirInput, "source-engine-*")
	if err != nil {
		logger.errorMsg("Failed to create secure temporary repository directory: " + err.Error())
		os.Exit(1)
	}
	Config.tempRepoDir = tempDir

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
		logger.errorMsg("Unsupported game: " + *gameBuildInput + ". Supported values are portal (not portal 2) and hl2 (half life 2).")
		os.Exit(1)
	}

	if *testStuff {
		logLevel = 3
		logger.debugMsg("Debug")
		logger.infoMsg("Info")
		logger.warnMsg("Warn")
		logger.successMsg("Success")
		logger.errorMsg("Error D:")

		Config.skipCleanup = true

	}

	if !checkSteamBetaRequirement(Config.GameToBuild) {
		os.Exit(1)
	}

	success := build()

	if success && !Config.dryRun {
		logger.infoMsg("Game has been successfully built and installed (I hope, go test it ig)!")
		logger.infoMsg("Would you like to delete this builder tool now? (If it works just do yes it takes up space) (y/n)")
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

	if strings.Contains(executable, ".app/Contents/MacOS/") {
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

	// fallback
	logger.infoMsg("Deleting executable: " + executable)
	err = os.Remove(executable)
	if err != nil {
		logger.errorMsg("Failed to delete executable: " + err.Error())
		// fallback to rm just in case
		exec.Command("rm", executable).Start()
	}
	logger.successMsg("Cleanup initiated. Goodbye and thanks for using my tool!")
}
