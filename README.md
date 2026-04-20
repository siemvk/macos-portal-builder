# Macos Source Builder Cli

This is a simple CLI tool to build portal and hl2 for macOS.

**MAKE SURE TO INSTALL THE STEAM BETA OR LEGACY VERSION OF THE GAME!!!**

**TO USE THIS SOFTWARE YOU MUST HAVE A LEGAL COPY OF THE GAME YOU WANT TO BUILD. PLEASE DO NOT ASK ME FOR HELP IF YOU DO NOT HAVE A LEGAL COPY OF THE GAME.**

## demo

### part 1: building the game

https://youtu.be/o51p2zmxCSo

### part 2: running the game

https://youtu.be/TXvaL4L8V6s

## Why

Valve has not updated the official versions of hl2 or portal to run on 64 bit macOS versions. But they did make the source engine source available and others have fixed it to run on modern macOS versions! But the process of building the game was annoying and the guide was not very clear. So I made this tool to make it easier for people to build the games because I really like these games and I want to be able to play them on my mac without having to use a old mac or a virtual machine. This was made even worse by the fact that steam dropped support for the last macos version that could run 32 bit apps, so now you have to use a windows or linux machine to play the game or build it yourself. I hope this tool will make it easier for people to play these games on their macs and maybe even get more people to play these games!

## How to use

1. download the steam beta or legacy version of the game you want to build
2. download the app from the releases tab, unzip and run it.
3. enter the name of the game you want to build (portal or hl2)
4. wait
5. delete the app because you don't need it anymore and it takes up space
6. launch the game from steam and enjoy!

## does it run any good?

Yes! On my m1 macbook air 2020 8gb ram it run at 300 fps on portal and I don't have any performance issues. I haven't tested hl2 yet but I expect it to run a bit worse.

## credits

This is only a UI built on top of the wonderful work of the people who have fixed the source engine to run on modern macOS versions. You can find their work here:

https://github.com/nillerusr/source-engine

And the guide I used to build the game is here:

https://jxhug.notion.site/Guide-to-Installing-Portal-Using-Source-Engine-on-macOS-660803f9ced149cfa1647d38fd5a7092

## ai usage

Ai was my Go emotional support line. This is because I am a complete beginner at Go.

## Whats up with the scrapped GUI?

The GUI was scrapped because I don't know how to make a GUI in Go and the CLI was easier to make. 
