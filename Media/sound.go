package main

// Original, deterministic sound design. Layered transients, pitch envelopes,
// filtered noise and short stereo reflections; no external samples or services.
import (
 "encoding/binary"
 "fmt"
 "math"
 "os"
 "path/filepath"
)

const rate = 48000
var names = []string{"shot", "impact", "pickup", "power", "bounce", "move", "rotate", "lock", "clear", "start", "win", "lose", "pause", "select"}
type sample struct { pcm []float64 }
func design(name string) sample {
 duration:=0.32
 switch name {case "shot":duration=.19;case "move":duration=.075;case "bounce":duration=.15;case "rotate":duration=.14;case "pickup":duration=.13;case "power", "clear":duration=.75;case "win":duration=1.4;case "lose":duration=.85;case "start":duration=.65}
 mono:=make([]float64,int(duration*rate)); phase:=0.0; noise:=0.0; state:=uint32(91273)
 for i:=range mono {
  t:=float64(i)/rate; u:=t/duration
  state=1664525*state+1013904223; white:=float64(state)/2147483648-1
  noise += .16*(white-noise)
  f:=440.0; amp:=math.Exp(-7*u); v:=0.0
  switch name {
  case "shot": f=1800*math.Exp(-17*t)+110; amp=math.Exp(-21*t);v=.15*white*math.Exp(-70*t)
  case "impact": f=105*math.Exp(-9*t)+34; amp=math.Exp(-13*t);v=1.2*noise*math.Exp(-14*t)+.18*white*math.Exp(-55*t)
  case "pickup": f=920+1100*t; amp=math.Exp(-23*t)
  case "power", "clear", "win", "start":
   notes:=[]float64{261.63,329.63,392,523.25,659.25,783.99}
   step:=int(t/.09); if step>=len(notes){step=len(notes)-1};f=notes[step]
   amp=math.Min(1,t/.012)*math.Exp(-2.8*u);v=.23*math.Sin(2*math.Pi*f*1.5*t)*amp
  case "bounce": f=680*math.Exp(-5*t);amp=math.Exp(-24*t)
  case "move": f=240; amp=math.Exp(-50*t)
  case "rotate":f=390+2700*t;amp=math.Exp(-25*t)
  case "lock":f=90+90*math.Exp(-35*t);amp=math.Exp(-18*t);v=.4*noise*math.Exp(-27*t)
  case "lose": f=260*math.Exp(-2.5*t)+40;amp=math.Exp(-3.4*t);v=.24*noise*amp
  case "pause":f=520-400*t
  case "select":f=660+660*t
  }
  phase+=2*math.Pi*f/rate
  v+=(math.Sin(phase)+.19*math.Sin(2*phase)+.08*math.Sin(3*phase))*amp
  // Click-free attack/release, enough headroom for simultaneous voices.
  edge:=math.Min(1,t/.003)*math.Min(1,(duration-t)/.018)
  mono[i]=v*edge*.21
 }
 stereo:=make([]float64,len(mono)*2)
 for i,v:=range mono {
  l,r:=v,v
  if i>int(.037*rate){l+=mono[i-int(.037*rate)]*.16}
  if i>int(.053*rate){r+=mono[i-int(.053*rate)]*.14}
  stereo[2*i]=l;stereo[2*i+1]=r
 }
 return sample{stereo}
}
func wav(path string, s sample) error {
 data:=make([]byte,44+len(s.pcm)*2);copy(data,"RIFF");binary.LittleEndian.PutUint32(data[4:],uint32(len(data)-8));copy(data[8:],"WAVEfmt ")
 binary.LittleEndian.PutUint32(data[16:],16);binary.LittleEndian.PutUint16(data[20:],1);binary.LittleEndian.PutUint16(data[22:],2)
 binary.LittleEndian.PutUint32(data[24:],rate);binary.LittleEndian.PutUint32(data[28:],rate*4);binary.LittleEndian.PutUint16(data[32:],4);binary.LittleEndian.PutUint16(data[34:],16)
 copy(data[36:],"data");binary.LittleEndian.PutUint32(data[40:],uint32(len(s.pcm)*2))
 for i,v:=range s.pcm {binary.LittleEndian.PutUint16(data[44+i*2:],uint16(int16(math.Max(-1,math.Min(1,v))*32767)))}
 return os.WriteFile(path,data,0644)
}
func render(dir string) error {
 if err:=os.MkdirAll(dir,0755);err!=nil{return err}
 var demo []float64
 for _,name:=range names {
  s:=design(name);if err:=wav(filepath.Join(dir,name+".wav"),s);err!=nil{return err}
  demo=append(demo,s.pcm...);demo=append(demo,make([]float64,rate/2)...)
 }
 if err:=wav(filepath.Join(dir,"showcase.wav"),sample{demo});err!=nil{return err}
 fmt.Println("Rendered",len(names),"original stereo effects and showcase.wav");return nil
}
