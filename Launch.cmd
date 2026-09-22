@echo off
setlocal
cd /d "%~dp0"
if exist "%~dp0Arcade.exe" set "ARCADE_BR_EXE=%~dp0Arcade.exe"
if defined ARCADE_BR_EXE goto configured
if exist "%~dp0runtime\Arcade.exe" set "ARCADE_BR_EXE=%~dp0runtime\Arcade.exe"
:configured
if not defined ARCADE_BR_EXE goto missing
if not exist "%ARCADE_BR_EXE%" goto missing
if not exist "%~dp0config\%~1.sys" goto badgame
start "%~1" "%ARCADE_BR_EXE%" -"%~dp0config\%~1.sys"
exit /b 0
:missing
echo The Arcade needs an installed, licensed Business Rules! runtime.
echo Put Arcade.exe beside Launch.cmd, or set ARCADE_BR_EXE to your runtime.
echo See README.md for setup. Runtime and license files are not included.
pause
exit /b 1
:badgame
echo Unknown Arcade launcher: %~1
exit /b 1
