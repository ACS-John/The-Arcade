// Speak - reads a text file and reads it aloud via the Windows OneCore
// speech engine (through PowerShell's Windows.Media.SpeechSynthesis WinRT
// API), for Aphelion.panda's optional story narration
// (Aphelion/Aphelion.panda).
//
// OneCore rather than classic desktop SAPI (System.Speech) because OneCore
// is where Windows registers voices added via Settings > Time & Language >
// Speech > Manage voices - classic SAPI only ever sees the small fixed set
// baked into desktop Windows (David/Zira). WinRT speech synthesis has no
// direct "speak now" call: it renders to an in-memory stream first, which
// is then played back with System.Media.SoundPlayer.
//
// BR launches this with SYSTEM -C (see Aphelion.panda's Speak: subroutine),
// which does not wait for it - so this program is free to block for as long
// as the narration takes without ever stalling the game. The text itself is
// never passed on a command line (prose can contain quotes, and arbitrarily
// long command lines are a real Windows limit); instead both this program
// and the PowerShell it launches read it from the same file.
//
// Usage:
//
//	Speak.exe <path-to-text-file> [-rate N] [-voice NAME]
//
//	-rate N      speaking rate, -10 (slowest) to 10 (fastest); default 0.
//	             Mapped onto WinRT's SpeakingRate (0.5-1.5, 1.0 normal) -
//	             an approximation, not the same scale SAPI used.
//	-voice NAME  installed voice's DisplayName (e.g. "Microsoft Mark");
//	             default is the system's default voice. List what's
//	             installed with (in PowerShell):
//	               [Windows.Media.SpeechSynthesis.SpeechSynthesizer,Windows.Media.SpeechSynthesis,ContentType=WindowsRuntime] | Out-Null
//	               [Windows.Media.SpeechSynthesis.SpeechSynthesizer]::AllVoices
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Accept options on either side of the filename, as used by the game and talk.
// flag.Parse alone stops at the first filename and silently ignores later options.
func parseArgs(args []string) (textFile string, rate int, voice string, err error) {
	flags := flag.NewFlagSet("Speak", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.IntVar(&rate, "rate", 0, "speaking rate, -10 to 10")
	flags.StringVar(&voice, "voice", "", "installed voice's DisplayName")
	var options, files []string
	for idx := 0; idx < len(args); idx++ {
		arg := args[idx]
		if arg == "--" {
			files = append(files, args[idx+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			options = append(options, arg)
			if arg == "-rate" || arg == "-voice" || arg == "--rate" || arg == "--voice" {
				idx++
				if idx >= len(args) {
					return "", 0, "", fmt.Errorf("missing value for %s", arg)
				}
				options = append(options, args[idx])
			}
		} else {
			files = append(files, arg)
		}
	}
	if err = flags.Parse(options); err != nil {
		return
	}
	if len(files) != 1 {
		err = fmt.Errorf("exactly one text file path is required")
		return
	}
	if rate < -10 || rate > 10 {
		err = fmt.Errorf("rate must be between -10 and 10")
		return
	}
	textFile = files[0]
	return
}

func main() {
	textFile, rate, voice, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "speak: %v\n", err)
		os.Exit(2)
	}
	if _, err := os.Stat(textFile); err != nil {
		fmt.Fprintf(os.Stderr, "speak: %v\n", err)
		os.Exit(1)
	}

	// -10..10 -> 0.5..1.5, matching WinRT SpeechSynthesisOptions.SpeakingRate's
	// centering on 1.0 for normal speed.
	speakingRate := 1.0 + float64(rate)*0.05

	// Double any single quote so it can't break out of the PowerShell
	// string literals below (PowerShell's own escape for ' inside a
	// single-quoted string is '' - not a backslash).
	psSafePath := strings.ReplaceAll(textFile, "'", "''")
	selectVoice := ""
	if voice != "" {
		psSafeVoice := strings.ReplaceAll(voice, "'", "''")
		selectVoice = fmt.Sprintf(
			`$v = [Windows.Media.SpeechSynthesis.SpeechSynthesizer]::AllVoices | Where-Object { $_.DisplayName -eq '%s' }; `+
				`if ($v) { $s.Voice = $v }; `,
			psSafeVoice,
		)
	}
	script := fmt.Sprintf(
		`$ErrorActionPreference = 'Stop'; Add-Type -AssemblyName System.Runtime.WindowsRuntime; `+
			`[Windows.Media.SpeechSynthesis.SpeechSynthesizer,Windows.Media.SpeechSynthesis,ContentType=WindowsRuntime] | Out-Null; `+
			`Function Await($t,$rt) { `+
			`$m = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object { $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -like 'IAsyncOperation*' })[0]; `+
			`$task = $m.MakeGenericMethod($rt).Invoke($null,@($t)); $task.Wait(-1) | Out-Null; $task.Result }; `+
			`$s = New-Object Windows.Media.SpeechSynthesis.SpeechSynthesizer; `+
			`$s.Options.SpeakingRate = %v; `+
			`%s`+
			`$text = [IO.File]::ReadAllText('%s'); `+
			`$stream = Await ($s.SynthesizeTextToStreamAsync($text)) ([Windows.Media.SpeechSynthesis.SpeechSynthesisStream]); `+
			`$netStream = [System.IO.WindowsRuntimeStreamExtensions]::AsStreamForRead($stream); `+
			`$player = New-Object System.Media.SoundPlayer; `+
			`$player.Stream = $netStream; `+
			`$player.PlaySync()`,
		speakingRate, selectVoice, psSafePath,
	)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "speak: speech synthesis or playback failed: %v; run in the signed-in Windows user's session with access to speech and audio devices\n", err)
		os.Exit(1)
	}
}
