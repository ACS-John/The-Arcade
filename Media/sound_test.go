package main
import("math";"testing";"encoding/binary";"os";"path/filepath")
func TestEffects(t *testing.T){
 for _,n:=range names{t.Run(n,func(t *testing.T){
  s:=design(n);peak:=0.0;energy:=0.0
  if len(s.pcm)%2!=0||len(s.pcm)<rate/10{t.Fatal("invalid stereo duration")}
  for _,v:=range s.pcm{if math.IsNaN(v)||math.IsInf(v,0){t.Fatal("nonfinite sample")};peak=math.Max(peak,math.Abs(v));energy+=v*v}
  if peak>.9||peak<.03||energy<1{t.Fatalf("bad level peak=%v energy=%v",peak,energy)}
  if math.Abs(s.pcm[0])>.001||math.Abs(s.pcm[len(s.pcm)-1])>.01{t.Fatal("click at endpoint")}
  p:=filepath.Join(t.TempDir(),"test.wav");if err:=wav(p,s);err!=nil{t.Fatal(err)};b,_:=os.ReadFile(p)
  if string(b[:4])!="RIFF"||int(binary.LittleEndian.Uint32(b[40:]))!=len(s.pcm)*2{t.Fatal("invalid WAV")}
 })}
}
