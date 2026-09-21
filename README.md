# The Arcade

Eight games in a custom neon 1980s arcade cabinet. Seven games run in Business Rules! (BR); Qix uses a native Windows Go engine for animated vector graphics.

Start **Arcade.cmd**, then click a game or type its number and press Enter. Every game also has its own .cmd shortcut. Qix.cmd runs independently without BR.

## Getting started

The Windows game binaries, original artwork, and sound effects are included. Clone or download this repository into a folder such as **C:\ACS\The Arcade**.

The seven BR games and the Arcade menu require a separately installed, licensed Business Rules! runtime. That proprietary runtime and its license are not included and are not covered by this project's Unlicense. Set the ARCADE_BR_EXE environment variable to your runtime executable. The launchers also recognize runtime\br.exe or an adjacent Dev-5\ACS 5.exe installation. No accounting application is started.

Qix needs Windows and the included Qix\Qix.exe and Media\ArcadeAudio.exe. Playing does not require Go, Node.js, downloads, or an account.

## Games and controls

| Game | Play | Controls |
| --- | --- | --- |
| BR Invaders | Formations, shields, enemy fire, bonus saucers | Left/Right or A/D; Space fires |
| Paddle Panic | Brick walls, angled rebounds, wide paddle, multiball | Left/Right or A/D; Space serves |
| Byte Muncher | Maze chasing, power pellets, four hunters | Arrows or WASD |
| Crosswalk Chaos | Traffic, logs, five homes, countdown | Arrows or WASD |
| Bitipede | Splitting swarm, mushrooms, spiders | Arrows or WASD; Space fires |
| Bitris | Falling blocks, hold, landing guide, line clears | [Full controls](Bitris/README.md) |
| Aphelion | Illustrated alien-world adventure with narration | Click story choices |
| Qix | Close regions to claim 75% while avoiding Qix and sparks | Arrows or WASD automatically move and draw |

Action games use **Enter** to start/retry, **P** to pause, **M** to mute effects, and **Q/Esc** to quit. In the menu, type **M + Enter**. Aphelion has separate effects and narration controls. Games launched from the menu return there on exit.

Qix starts with two moving points on stage one, three on stage two, and four from stage three onward. Drawing begins automatically when you leave a boundary at normal speed. Hold Shift for slow drawing; a line drawn entirely in slow mode earns double points. Close your line against a safe boundary to claim territory. Crossing your own trail, letting Qix touch it, or meeting a spark costs a life. A fuse pursues you if you stop mid-line. Losing focus pauses play.

Qix saves best score and mute preference under the Windows user configuration directory in The Arcade\qix.json. The five BR arcade action games keep best scores in temporary files; cleanup can erase them. Aphelion and Bitris use their existing save files.

## Build and test

BR source is authored in .panda files, compiled with the ACS Nine-Tailed Fox compiler. Never reconstruct source by decompiling the included .br binaries.

Run Build.ps1 -Test to compile the menu, seven BR games, and six rule suites, run those suites, and smoke-test six rendered game startups. Supply -Compiler, -BrExe, and optionally -BrlsExe, or set ARCADE_COMPILER, ARCADE_BR_EXE, and ARCADE_BRLS. Adjacent Dev-5 tooling is recognized automatically. Omit -Test for compilation only. BR execution requires your installed runtime and license.

For Qix, install Go and run Qix\Build.ps1. It runs engine tests before building the native Windows executable. Tests cover automatic drawing, capture scoring, protected regions, trail collisions, the fuse, pause, stage progression, Qix complexity, and extended simulated play.

Media\Build.ps1 rebuilds the Go sound mixer, synthesized WAV effects, and PNG sprites. It requires Go plus Node.js with sharp; NodeExe and ModulesPath parameters can identify an existing installation. The sprite sources are original SVG files. Go components use only the standard library.

Tests\Native.ps1 and Tests\Capture.ps1 target a selected game process for UI checks. Generated checks, intermediate .brs files, local runtime/license files, and screenshots are ignored by Git.

## Graphics and sound

Six BR action games have individually illustrated backdrops and native PNG sprites. Qix combines an illustrated cabinet with glowing animated vectors and complete off-screen frame composition to prevent flicker. Aphelion retains its illustrated scenes and Windows speech narration.

Fourteen original stereo effects are synthesized from layered waveforms and filtered noise, with overlapping voices and soft limiting. Listen to [the sound showcase](audio/showcase.wav). The helper uses the default Windows audio device; games remain playable without it.

The menu and scene art were generated with the built-in image generation tool. Prompts are preserved in [art/prompt.txt](art/prompt.txt), [art/game-art-prompts.json](art/game-art-prompts.json), and [Qix/art-prompt.txt](Qix/art-prompt.txt). Editable sprite artwork and PNG exports live in art/sprites; Media/build-sprites.cjs reproduces them. Audio synthesis lives in Media/sound.go.

## License

This project's original code and assets are released under [The Unlicense](LICENSE). These are independently implemented arcade homages, without original game ROMs, recordings, or extracted artwork. Referenced game names remain associated with their respective owners; no affiliation is implied. The separately supplied BR runtime has its own license.