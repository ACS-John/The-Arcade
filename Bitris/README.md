# Bitris

A native Business Rules! falling-block game. Run **Bitris.cmd** in the Arcade folder (the root shortcut also works), or choose Bitris from **Developer → The Arcade**. Press Enter to start or retry after game over.

- Left/Right or A/D: move. Up/W/X: clockwise rotation. Z: counterclockwise.
- Down/S: soft drop. Space: hard drop. C: hold once per piece.
- P: pause/resume. Q/Esc: quit.

Seven-bag randomization, next-piece preview, landing guide, wall/floor kicks, a 450 ms lock delay (up to 15 grounded adjustments), and twenty visible board rows. Clearing 1/2/3/4 rows awards 100/300/500/800 points times the current level. Soft/hard drops award 1/2 points per row. Every ten cleared rows raises the level and falling speed.

The best score is saved locally in `%TEMP%\Bitris_best.txt`; it can disappear when Windows temporary files are cleared. No accounting data or ACS libraries are used by the standalone game. Exiting from ACS returns through its normal `Proc R` path. Rotation uses simple wall/floor kicks, not tournament-standard SRS.

Source is unnumbered `.panda`. `Engine.panda` is included by both the game and the regression tests. From the Arcade repository root, run Build.ps1 -Test using your installed BR runtime and Nine-Tailed Fox compiler; see the root README for setup. The test report must end with zero failures. The suite exercises all shape rotations, 100 seven-bags, board boundaries, occupied cells, adjacent and separated clears, scoring/levels, ghost/hard drop, hold limits, wall kicks, blocked spawn, restart, lock delay, and the lock-reset cap. Run BR from the signed-in Windows account; the Codex sandbox account can leave the runtime at an inaccessible startup error dialog.

Verified 2026-09-21: real BR compilation and 920 rule assertions passed. Native window and keyboard checks cover movement, rotation, drops, hold, pause, game over, and retry. `.brs`, test binaries, and reports are generated artifacts; `Bitris.br` is the playable build.
