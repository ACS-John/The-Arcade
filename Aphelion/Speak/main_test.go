package main

import "testing"

func TestOptionsBeforeOrAfterFilename(t *testing.T) {
	for _, args := range [][]string{
		{`C:\a folder\speech.txt`, "-rate", "-3", "-voice", "Microsoft Mark"},
		{"-rate", "-3", "-voice", "Microsoft Mark", `C:\a folder\speech.txt`},
		{"-voice=Microsoft Mark", `C:\a folder\speech.txt`, "-rate=-3"},
	} {
		path, rate, voice, err := parseArgs(args)
		if err != nil || path != `C:\a folder\speech.txt` || rate != -3 || voice != "Microsoft Mark" {
			t.Fatalf("%q: got %q, %d, %q, %v", args, path, rate, voice, err)
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{}, {"a.txt", "b.txt"}, {"a.txt", "-rate"},
		{"a.txt", "-rate", "11"}, {"a.txt", "-rate", "-11"},
		{"a.txt", "-unknown"}, {"a.txt", "-rate", "fast"},
	} {
		if _, _, _, err := parseArgs(args); err == nil {
			t.Fatalf("accepted invalid arguments %q", args)
		}
	}
}
