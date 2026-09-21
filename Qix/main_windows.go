package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//go:embed background.png
var background []byte
var user = syscall.NewLazyDLL("user32.dll")
var gdi = syscall.NewLazyDLL("gdi32.dll")
var kernel = syscall.NewLazyDLL("kernel32.dll")

func u(name string, args ...uintptr) uintptr { r, _, _ := user.NewProc(name).Call(args...); return r }
func d(name string, args ...uintptr) uintptr { r, _, _ := gdi.NewProc(name).Call(args...); return r }
func wide(s string) *uint16                  { p, _ := syscall.UTF16PtrFromString(s); return p }

type rect struct{ L, T, R, B int32 }
type paint struct {
	DC           uintptr
	Erase        int32
	Rect         rect
	Restore, Inc int32
	Reserved     [32]byte
}
type msg struct {
	H       uintptr
	Message uint32
	W, L    uintptr
	Time    uint32
	X, Y    int32
	Private uint32
}
type windowClass struct {
	Size, Style                        uint32
	Proc                               uintptr
	ClsExtra, WndExtra                 int32
	Instance, Icon, Cursor, Background uintptr
	Menu, Name                         *uint16
	Small                              uintptr
}
type bitmapInfo struct {
	Size                   uint32
	Width, Height          int32
	Planes, Bits           uint16
	Compression, ImageSize uint32
	XPels, YPels           int32
	Used, Important        uint32
}
type preferences struct {
	Best  int
	Muted bool
}

var state = newGame()
var screen *canvas
var hwnd uintptr
var last = time.Now()
var fonts = map[int]uintptr{}
var mailbox string
var audioCmd *exec.Cmd
var audioSeq int
var audioLast time.Time
var audioAvailable bool
var savePath string
var pausedByFocus bool

func loadSave() {
	base, e := os.UserConfigDir()
	if e != nil {
		return
	}
	savePath = filepath.Join(base, "The Arcade", "qix.json")
	b, e := os.ReadFile(savePath)
	if e == nil {
		var p preferences
		if json.Unmarshal(b, &p) == nil {
			state.Best = max(0, p.Best)
			state.Muted = p.Muted
		}
	}
}
func save() {
	if savePath == "" {
		return
	}
	os.MkdirAll(filepath.Dir(savePath), 0700)
	b, _ := json.Marshal(preferences{state.Best, state.Muted})
	tmp := savePath + ".tmp"
	if os.WriteFile(tmp, b, 0600) == nil {
		os.Rename(tmp, savePath)
	}
}
func audioStart() {
	exe, _ := os.Executable()
	helper := filepath.Join(filepath.Dir(exe), "..", "Media", "ArcadeAudio.exe")
	f, e := os.CreateTemp("", "TheArcade-Qix-*.txt")
	if e != nil {
		return
	}
	mailbox = f.Name()
	f.Close()
	audioAvailable = true
	audioFlush(false)
	audioCmd = exec.Command(helper, "-mailbox", mailbox)
	audioCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if e = audioCmd.Start(); e != nil {
		audioAvailable = false
		return
	}
}
func audioFlush(stop bool) {
	if !audioAvailable {
		return
	}
	if !stop && len(state.Events) == 0 && time.Since(audioLast) < 500*time.Millisecond {
		return
	}
	mode := "play"
	if state.Muted {
		mode = "mute"
	}
	if stop {
		mode = "stop"
	}
	audioSeq++
	os.WriteFile(mailbox, []byte(fmt.Sprintf("%d|%s|%s\n", audioSeq, mode, strings.Join(state.Events, ","))), 0600)
	audioLast = time.Now()
	state.Events = nil
}
func text(dc uintptr, x, y, size int, s string, col uint32, scale float64, ox, oy int) {
	height := max(10, int(float64(size)*scale))
	font := fonts[height]
	if font == 0 {
		font = d("CreateFontW", uintptr(-int32(height)), 0, 0, 0, 500, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(wide("Segoe UI"))))
		fonts[height] = font
	}
	old := d("SelectObject", dc, font)
	d("SetBkMode", dc, 1)
	d("SetTextColor", dc, uintptr((col&255)<<16|(col&0xff00)|(col>>16)))
	utf := syscall.StringToUTF16(s)
	d("TextOutW", dc, uintptr(ox+int(float64(x)*scale)), uintptr(oy+int(float64(y)*scale)), uintptr(unsafe.Pointer(&utf[0])), uintptr(len(utf)-1))
	d("SelectObject", dc, old)
}
func drawScene(dc uintptr) {
	var area rect
	u("GetClientRect", hwnd, uintptr(unsafe.Pointer(&area)))
	scale := math.Min(float64(area.R)/SW, float64(area.B)/SH)
	ww, hh := int(SW*scale), int(SH*scale)
	ox, oy := (int(area.R)-ww)/2, (int(area.B)-hh)/2
	d("PatBlt", dc, 0, 0, uintptr(area.R), uintptr(area.B), 0x42)
	screen.draw(state)
	info := bitmapInfo{Size: 40, Width: SW, Height: -SH, Planes: 1, Bits: 32}
	d("SetStretchBltMode", dc, 4)
	d("StretchDIBits", dc, uintptr(ox), uintptr(oy), uintptr(ww), uintptr(hh), 0, 0, SW, SH, uintptr(unsafe.Pointer(&screen.Pix[0])), uintptr(unsafe.Pointer(&info)), 0, 0xcc0020)
	t := func(x, y, size int, s string, c uint32) { text(dc, x, y, size, s, c, scale, ox, oy) }
	t(906, 188, 15, "TERRITORY / CONTROL", 0x81a6c0)
	t(906, 218, 38, fmt.Sprintf("%06d", state.Score), 0xf1f7ff)
	t(906, 270, 16, fmt.Sprintf("BEST  %06d", state.Best), 0x83a2b9)
	t(906, 315, 32, fmt.Sprintf("%.1f%%", state.Percent), 0x69eeee)
	t(1070, 330, 16, "GOAL 75%", 0xa8bad1)
	t(906, 401, 21, fmt.Sprintf("STAGE %02d    LIVES %d", state.Level, state.Lives), 0xedf2ff)
	t(906, 457, 17, "ARROWS / WASD    Move + draw", 0xc7d8eb)
	t(906, 490, 17, "Drawing starts automatically.", 0x65e9f2)
	t(906, 523, 17, "Hold SHIFT   Slow draw / 2x", 0xff9fc7)
	t(906, 573, 17, "P   Pause        ESC   Exit", 0xc7d8eb)
	sound := "M   Sound on"
	if state.Muted {
		sound = "M   Sound off"
	}
	if !audioAvailable {
		sound = "Sound unavailable"
	}
	t(906, 611, 17, sound, 0xc7d8eb)
	t(906, 680, 15, "KEEP MOVING. THE FUSE FOLLOWS.", 0xe5b989)
	t(65, 777, 20, "DRAW A LINE. CLOSE A REGION. CONTAIN THE QIX.", 0xaad8e8)
	t(65, 813, 15, "Hold Shift throughout a line for double points. Capture 75% to advance.", 0x90aabe)
	title, sub, action := "", "", ""
	switch state.Mode {
	case "ready":
		title = "CONTAIN THE CHAOS"
		sub = "Use arrows or WASD to move and draw."
		action = "ENTER TO PLAY"
	case "paused":
		title = "PAUSED"
		sub = "Your territory is waiting."
		action = "P TO RESUME"
	case "over":
		title = "SIGNAL LOST"
		sub = fmt.Sprintf("FINAL SCORE  %06d", state.Score)
		action = "ENTER TO TRY AGAIN"
	case "cleared":
		title = "SECTOR SECURED"
		sub = fmt.Sprintf("%.1f%% captured. Preparing the next stage.", state.Percent)
		action = "WELL PLAYED"
	}
	if title != "" {
		t(231, 391, 30, title, 0xeeefff)
		t(231, 439, 18, sub, 0xaccaE0)
		t(231, 478, 18, action, 0x74f7ef)
	}
}
func callback(h uintptr, m uint32, w, l uintptr) uintptr {
	switch m {
	case 0x14:
		return 1
	case 5:
		u("InvalidateRect", h, 0, 0)
		return 0
	case 6:
		if w&0xffff == 0 {
			state.Keys = [256]bool{}
			if state.Mode == "playing" {
				state.Mode = "paused"
				pausedByFocus = true
			}
		}
		return 0
	case 0x100, 0x104:
		if w < 256 {
			state.Keys[w] = true
		}
		if l&(1<<30) == 0 {
			switch w {
			case 13:
				if state.Mode == "ready" || state.Mode == "over" {
					state.restart()
				}
			case 'P':
				if state.Mode == "playing" {
					state.Mode = "paused"
					state.event("pause")
				} else if state.Mode == "paused" {
					state.Mode = "playing"
					pausedByFocus = false
				}
			case 'M':
				state.Muted = !state.Muted
				audioLast = time.Time{}
				save()
			case 27, 'Q':
				u("DestroyWindow", h)
			}
		}
		return 0
	case 0x101, 0x105:
		if w < 256 {
			state.Keys[w] = false
		}
		return 0
	case 0x113:
		now := time.Now()
		dt := math.Min(.05, now.Sub(last).Seconds())
		last = now
		for dt > 0 {
			step := math.Min(dt, 1.0/120)
			state.tick(step)
			dt -= step
		}
		audioFlush(false)
		u("InvalidateRect", h, 0, 0)
		return 0
	case 0x318:
		draw(w)
		return 0
	case 0xf:
		var ps paint
		dc := u("BeginPaint", h, uintptr(unsafe.Pointer(&ps)))
		draw(dc)
		u("EndPaint", h, uintptr(unsafe.Pointer(&ps)))
		return 0
	case 2:
		save()
		audioFlush(true)
		for _, f := range fonts {
			d("DeleteObject", f)
		}
		u("PostQuitMessage", 0)
		return 0
	}
	return u("DefWindowProcW", h, uintptr(m), w, l)
}
func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	bg, e := png.Decode(bytes.NewReader(background))
	if e != nil {
		panic(e)
	}
	screen = newCanvas(bg)
	loadSave()
	if len(os.Args) > 1 && os.Args[1] == "-smoke" {
		state.restart()
		state.Player = point{25, 0}
		state.Anchor = state.Player
		for y := 1; y < H; y++ {
			state.move(1)
		}
		for i := 0; i < 180; i++ {
			state.tick(1.0 / 120)
		}
		screen.draw(state)
		if len(os.Args) > 2 {
			if e = screen.save(os.Args[2]); e != nil {
				panic(e)
			}
		}
		fmt.Printf("Qix smoke: %.1f%%, score %d, lives %d\n", state.Percent, state.Score, state.Lives)
		return
	}
	instance, _, _ := kernel.NewProc("GetModuleHandleW").Call(0)
	cls := wide("TheArcadeQix")
	wc := windowClass{Size: uint32(unsafe.Sizeof(windowClass{})), Style: 3, Proc: syscall.NewCallback(callback), Instance: instance, Cursor: u("LoadCursorW", 0, 32512), Name: cls}
	if u("RegisterClassExW", uintptr(unsafe.Pointer(&wc))) == 0 {
		panic("RegisterClassEx failed")
	}
	hwnd = u("CreateWindowExW", 0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(wide("Qix | The Arcade"))), 0x00cf0000, 80, 50, 1280, 900, 0, 0, instance, 0)
	if hwnd == 0 {
		panic("CreateWindowEx failed")
	}
	audioStart()
	u("ShowWindow", hwnd, 5)
	u("UpdateWindow", hwnd)
	u("SetTimer", hwnd, 1, 16, 0)
	last = time.Now()
	var message msg
	for u("GetMessageW", uintptr(unsafe.Pointer(&message)), 0, 0, 0) != 0 {
		u("TranslateMessage", uintptr(unsafe.Pointer(&message)))
		u("DispatchMessageW", uintptr(unsafe.Pointer(&message)))
	}
	if audioCmd != nil && audioCmd.Process != nil {
		audioCmd.Wait()
	}
	if mailbox != "" {
		os.Remove(mailbox)
	}
}

var bufferDC, bufferBitmap, bufferOld uintptr
var bufferW, bufferH int32

func draw(target uintptr) {
	var area rect
	u("GetClientRect", hwnd, uintptr(unsafe.Pointer(&area)))
	if area.R <= 0 || area.B <= 0 {
		return
	}
	if bufferDC == 0 || bufferW != area.R || bufferH != area.B {
		if bufferDC != 0 {
			d("SelectObject", bufferDC, bufferOld)
			d("DeleteObject", bufferBitmap)
			d("DeleteDC", bufferDC)
		}
		bufferDC = d("CreateCompatibleDC", target)
		bufferBitmap = d("CreateCompatibleBitmap", target, uintptr(area.R), uintptr(area.B))
		bufferOld = d("SelectObject", bufferDC, bufferBitmap)
		bufferW = area.R
		bufferH = area.B
	}
	drawScene(bufferDC)
	d("BitBlt", target, 0, 0, uintptr(area.R), uintptr(area.B), bufferDC, 0, 0, 0xcc0020)
}
