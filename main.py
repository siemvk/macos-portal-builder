import os
from CLogging import logging, LogLevel
import argparse
import sys
import enum
import tkinter as tk

logger = logging("app.log", LogLevel.DEBUG)


class Game(enum.Enum):
    PORTAL = "portal"
    HALFLIFE2 = "hl2"


class BuildConfigClass:
    def __init__(self):
        pass

    repoURL = "https://github.com/nillerusr/source-engine"
    tempRepoDir = os.path.expanduser("~/.temp/source-engine-build-tool")
    gameToBuild: Game = Game.PORTAL
    skipCleanup: bool = False

    showCommandOutput: bool = False
    dryRun: bool = False
    skipBuild: bool = False


buildConfig = BuildConfigClass()


def exec(command):
    logger.debug(f"Executing command: {command}")
    if not buildConfig.showCommandOutput:
        return os.system(command + " >& /dev/null")
    else:
        return os.system(command)


def portal_only_res_fix():
    import subprocess
    import re
    import os

    cfg_path = os.path.expanduser(
        "~/Library/Application Support/Steam/steamapps/common/Portal/portal/cfg/autoexec.cfg"
    )

    try:
        output = subprocess.check_output(
            ["system_profiler", "SPDisplaysDataType"]
        ).decode()

        match = re.search(r"Resolution: (\d+) x (\d+)", output)
        if not match:
            print("Could not detect resolution")
            return

        w, h = map(int, match.groups())

        # fallback for high (Retina) resolutions
        if w > 2560:
            w, h = w // 2, h // 2

        os.makedirs(os.path.dirname(cfg_path), exist_ok=True)

        with open(cfg_path, "w") as f:
            f.write(f"mat_setvideomode {w} {h} 1\n")
            f.write("mat_fullscreen 1\n")

        print(f"Set resolution to {w}x{h}")

    except Exception as e:
        print("Error:", e)


def clone_repo():
    if not os.path.exists(buildConfig.tempRepoDir):
        exec(f"git clone --recursive {buildConfig.repoURL} {buildConfig.tempRepoDir}")
    else:
        logger.error(
            f"Directory {buildConfig.tempRepoDir} already exists. Skipping clone."
        )
        logger.warn(f"This may cause issues if this is not the expected repo.")


def cleanup():
    if buildConfig.skipCleanup:
        logger.warn("Skipping cleanup as per configuration.")
        return
    if os.path.exists(buildConfig.tempRepoDir):
        exec(f"rm -rf {buildConfig.tempRepoDir}")
    else:
        logger.error(
            f"Directory {buildConfig.tempRepoDir} does not exist. Some weird shit is going on."
        )


def gameNameToDir(gameName: Game):
    if gameName == Game.PORTAL:
        path = os.path.expanduser(
            "~/Library/Application\ Support/Steam/steamapps/common/Portal"
        )
        os.makedirs(path, exist_ok=True)
        return path
    elif gameName == Game.HALFLIFE2:
        path = os.path.expanduser(
            "~/Library/Application\ Support/Steam/steamapps/common/Half-Life 2"
        )
        os.makedirs(path, exist_ok=True)
        return path


def CLI():

    logger.info(
        "Please make sure you have the older version (via steam) of portal/hl2 else this script wil not work!"
    )
    input("press enter to continue")

    parser = argparse.ArgumentParser(description="MacOS Source Builder CLI")
    parser.add_argument(
        "--gameName",
        help="Name of the game to build (for example: --gameName=portal)",
        required=True,
        type=Game,
    )
    parser.add_argument(
        "--overrideRepoURL",
        help="Override the default repository URL",
        required=False,
        type=str,
    )
    parser.add_argument(
        "--overrideTempRepoDir",
        help="Override the default temporary repository directory",
        required=False,
        type=str,
    )
    parser.add_argument(
        "--logLevel",
        help="Set the logging level (DEBUG, INFO, WARNING, ERROR) (NOTE: You can not see logs that are before the cli is loaded, so if you want to see more logs you have to fork and change the default log level in the code)",
        required=False,
        type=LogLevel,
    )
    parser.add_argument(
        "--skipCleanup",
        help="Skip the cleanup step (WARNING: This will leave the temporary repository directory on your system)",
        required=False,
        type=bool,
    )
    parser.add_argument(
        "--showCommandOutput",
        help="Show the output of the executed commands (WARNING: This will make the output very verbose)",
        required=False,
        type=bool,
    )
    parser.add_argument(
        "--skipInstall",
        help="dry run the build process without actually installing the game",
        required=False,
        type=bool,
    )
    parser.add_argument(
        "--cleanBeforeStarting",
        help="Clean the temporary repository directory before starting the build process",
        required=False,
        type=bool,
    )
    parser.add_argument(
        "--skipBuild",
        help="Skip the build step",
        required=False,
        type=bool,
    )
    args = parser.parse_args()
    buildConfig.repoURL = (
        args.overrideRepoURL if args.overrideRepoURL else buildConfig.repoURL
    )
    buildConfig.tempRepoDir = (
        args.overrideTempRepoDir
        if args.overrideTempRepoDir
        else buildConfig.tempRepoDir
    )
    logger.log_level = args.logLevel if args.logLevel else logger.log_level
    buildConfig.gameToBuild = args.gameName
    buildConfig.skipCleanup = (
        args.skipCleanup if args.skipCleanup else buildConfig.skipCleanup
    )
    buildConfig.showCommandOutput = (
        args.showCommandOutput
        if args.showCommandOutput
        else buildConfig.showCommandOutput
    )
    buildConfig.skipBuild = args.skipBuild if args.skipBuild else buildConfig.skipBuild
    buildConfig.dryRun = args.skipInstall if args.skipInstall else buildConfig.dryRun
    if buildConfig.skipCleanup:
        logger.warn(
            "Cleanup will be skipped. DO NOT COMPLAIN ABOUT LEFTOVER FILES IN YOUR SYSTEM."
        )
    if args.cleanBeforeStarting:
        logger.info("Cleaning temporary repository directory before starting...")
        cleanup()
        logger.success("done cleaning temporary repository directory!")
    logger.debug(f"CLI arguments parsed: {args}")

    build()


def build():
    logger.debug(f"Starting build process for game: {buildConfig.gameToBuild.value}")

    logger.info("Cloning repository...")
    clone_repo()
    logger.success("done cloning repository")

    logger.info("Installing dependencies...")
    logger.debug("Using Homebrew to install dependencies. This may take a while...")
    exec(
        "brew install sdl2 freetype2 fontconfig pkg-config opus jpeg jpeg-turbo libpng libedit"
    )
    logger.debug("Installing Xcode Command Line Tools. This may take a while...")
    exec("xcode-select --install")
    logger.success("done installing dependencies!")

    logger.info("Configureing build script...")
    try1 = exec(
        f"cd {buildConfig.tempRepoDir} && python3 waf configure -T release --prefix='' --build-games={buildConfig.gameToBuild.value}"
    )
    if try1 != 0:
        logger.error(
            "Basic install failed! this is not uncommon, trying again with different clang"
        )
        try2 = exec(
            f"cd {buildConfig.tempRepoDir} && export CC=/usr/bin/clang && export CXX=/usr/bin/clang++ && python3 waf configure -T release --prefix='' --build-games={buildConfig.gameToBuild.value}"
        )
        if try2 != 0:
            logger.error(
                "Install failed again! I do not experience this on my machine, so I am doing random fixes from reddit now."
            )
            try3 = exec(
                f"cd {buildConfig.tempRepoDir} && export CC=/usr/bin/clang && export CXX=/usr/bin/clang++ && arch -arm64 python3 waf configure -T release --prefix='' --build-games={buildConfig.gameToBuild.value}"
            )
            if try3 != 0:
                logger.error(
                    "Install failed again!!!! Oke so what if the first fix broke the second fix so lets try the second fix without the first fix."
                )
                try4 = exec(
                    f"cd {buildConfig.tempRepoDir} && arch -arm64 python3 waf configure -T release --prefix='' --build-games={buildConfig.gameToBuild.value}"
                )
                if try4 != 0:
                    logger.error(
                        "Install failed again!!!!! I give up. Please open an issue with the log output so I can try to fix this."
                    )
                    if not buildConfig.skipCleanup:
                        logger.info("Cleaning up temporary repository directory...")
                        cleanup()
                        logger.success(
                            "done cleaning up temporary repository directory!"
                        )
                    logger.error("Open a issue!!!!")
                    sys.exit(1)
    logger.success("done configuring build script!")

    if not buildConfig.skipBuild:
        logger.info("Building the game...")
        exec(f"cd {buildConfig.tempRepoDir} && python3 waf build")
        logger.success("done building the game!")
    else:
        logger.warn("Skipping build step as per configuration.")

    logger.info("Installing the game to a temp directory...")
    exec(
        f"cd {buildConfig.tempRepoDir} && python3 waf install --destdir={buildConfig.tempRepoDir}/installingthismf"
    )

    logger.success("done installing the game!")

    if buildConfig.dryRun:
        logger.warn("Dry run enabled, skipping installation to game folder.")
        cleanup()
        return
    logger.info("copying files to the game folder...")
    logger.debug(
        f"Copying files from {buildConfig.tempRepoDir}/installingthismf to {gameNameToDir(buildConfig.gameToBuild)}"
    )
    os.chdir(gameNameToDir(buildConfig.gameToBuild))
    logger.debug("Removing old bin files...")
    exec("rm -rf ./portal/bin ./bin")
    logger.debug("Copying new bin files...")
    exec(f"cp -r {buildConfig.tempRepoDir}/installingthismf/portal/bin ./portal/bin")
    exec(f"cp -r {buildConfig.tempRepoDir}/installingthismf/bin ./bin")
    logger.debug("renaming old hl2_osx file to hl2_osx_backup...")
    exec("mv ./hl2_osx ./hl2_osx_backup")
    logger.debug("copying new hl2_launcher as hl2_osx file...")
    exec(
        f"mv {buildConfig.tempRepoDir}/installingthismf/hl2_launcher {gameNameToDir(buildConfig.gameToBuild)}/hl2_osx"
    )

    logger.success("done copying files to the game folder!")
    logger.info("Just cleaning up...")
    cleanup()
    logger.success("done cleaning up")
    if buildConfig.gameToBuild == Game.PORTAL:
        logger.info("Applying Portal resolution fix...")
        portal_only_res_fix()
        logger.success("done applying Portal resolution fix!")


def GUI():
    root = tk.Tk()
    root.mainloop()
    root.title("Macos Source Builder")


def main():
    logger.debug("checking for CLI arguments to determine the mode of operation")
    if len(sys.argv) > 1:
        logger.info("CLI arguments detected, executing CLI()")
        CLI()
    else:
        logger.info("No CLI arguments detected, executing GUI()")
        GUI()


if __name__ == "__main__":
    logger.debug("Application starting...")
    main()
else:
    logger.debug("Module imported, not executing main()")
