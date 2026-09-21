package main

import (
 "flag"
 "fmt"
 "math"
 "os"
 "runtime"
 "strings"
 "syscall"
 "time"
 "unsafe"
)

var winmm=syscall.NewLazyDLL("winmm.dll")
var waveOpen=winmm.NewProc("waveOutOpen")
var wavePrepare=winmm.NewProc("waveOutPrepareHeader")
var waveWrite=winmm.NewProc("waveOutWrite")
var waveReset=winmm.NewProc("waveOutReset")
var waveUnprepare=winmm.NewProc("waveOutUnprepareHeader")
var waveClose=winmm.NewProc("waveOutClose")
type format struct { Tag,Channels uint16; Rate,Bytes uint32; Align,Bits,Extra uint16 }
type header struct { Data uintptr; Length,Recorded uint32; User uintptr; Flags,Loops uint32; Next,Reserved uintptr }
type voice struct { sound sample; at int }

func play(mailbox string) error {
 var device uintptr
 f:=format{Tag:1,Channels:2,Rate:rate,Bytes:rate*4,Align:4,Bits:16}
 code,_,_:=waveOpen.Call(uintptr(unsafe.Pointer(&device)),^uintptr(0),uintptr(unsafe.Pointer(&f)),0,0,0)
 if code!=0{return fmt.Errorf("waveOutOpen: %d",code)}
 const blocks=4; const frames=480
 buffers:=make([][]int16,blocks);headers:=make([]header,blocks)
 defer func(){waveReset.Call(device);for i:=range headers {waveUnprepare.Call(device,uintptr(unsafe.Pointer(&headers[i])),unsafe.Sizeof(header{}))};waveClose.Call(device);runtime.KeepAlive(buffers);runtime.KeepAlive(headers)}()
 bank:=map[string]sample{};for _,n:=range names{bank[n]=design(n)}
 for i:=range headers {
  buffers[i]=make([]int16,frames*2);headers[i]=header{Data:uintptr(unsafe.Pointer(&buffers[i][0])),Length:frames*4}
  if result,_,_:=wavePrepare.Call(device,uintptr(unsafe.Pointer(&headers[i])),unsafe.Sizeof(header{}));result!=0{return fmt.Errorf("prepare: %d",result)}
 }
 // READY is evidence of an opened device, not merely of a process launch.
 os.WriteFile(mailbox+".log",[]byte("READY stereo 48000Hz / 40ms queue\n"),0600)
 voices:=[]voice{};last:="";lastMessage:=time.Now();muted:=false;events:=0;muteChanges:=0;tick:=time.NewTicker(5*time.Millisecond);defer tick.Stop()
 defer func(){os.WriteFile(mailbox+".log",[]byte(fmt.Sprintf("STOPPED stereo 48000Hz; mixed events=%d; mute changes=%d\n",events,muteChanges)),0600)}()
 submitted:=make([]bool,blocks)
 for range tick.C {
  if time.Since(lastMessage)>4*time.Second{return nil} // Parent crash or abandonment.
  raw,err:=os.ReadFile(mailbox)
  if err==nil && len(raw)>0 && raw[len(raw)-1]=='\n' {
   message:=strings.TrimSpace(string(raw));parts:=strings.SplitN(message,"|",3)
   if len(parts)==3 && parts[0]!=last {
    last=parts[0];lastMessage=time.Now()
    if parts[1]=="stop"{return nil}
    nextMuted:=parts[1]=="mute";if nextMuted!=muted{muteChanges++};muted=nextMuted
    if muted {voices=nil}
    for _,event:=range strings.Split(parts[2],",") {
     if s,ok:=bank[event];ok && !muted {
      if len(voices)>=16{voices=voices[1:]};voices=append(voices,voice{sound:s});events++
     }
    }
   }
  }
  for i:=range headers {
   if submitted[i] && headers[i].Flags&1==0{continue}
   for k:=range buffers[i] {
    v:=0.0
    for j:=range voices {if voices[j].at<len(voices[j].sound.pcm){v+=voices[j].sound.pcm[voices[j].at];voices[j].at++}}
    buffers[i][k]=int16(math.Tanh(v*.72)*30000)
   }
   alive:=voices[:0];for _,v:=range voices{if v.at<len(v.sound.pcm){alive=append(alive,v)}};voices=alive
   if result,_,_:=waveWrite.Call(device,uintptr(unsafe.Pointer(&headers[i])),unsafe.Sizeof(header{}));result!=0{return fmt.Errorf("waveOutWrite: %d",result)}
   submitted[i]=true
  }
 }
 return nil
}
func main(){
 out:=flag.String("render","","render WAV assets");mail:=flag.String("mailbox","","BR session mailbox");flag.Parse()
 var err error
 if *out!=""{err=render(*out)}else if *mail!=""{err=play(*mail)}else{err=fmt.Errorf("use -render directory or -mailbox file")}
 if err!=nil {if *mail!=""{os.WriteFile(*mail+".log",[]byte("ERROR "+err.Error()+"\n"),0600)};fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
