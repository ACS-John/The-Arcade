param([Parameter(Mandatory=$true)][int]$GameProcessId,[string]$Text,[string]$ButtonTitle)
Add-Type -TypeDefinition @"
using System; using System.Text; using System.Runtime.InteropServices;
public class ArcadeNative {
 public delegate bool EnumProc(IntPtr h,IntPtr p);
 [StructLayout(LayoutKind.Sequential)] public struct RECT {public int l,t,r,b;}
 [StructLayout(LayoutKind.Sequential)] public struct GUI {public int cb; public uint flags; public IntPtr active,focus,capture,menu,move,caret;public RECT rect;}
 [DllImport("user32.dll")] public static extern uint GetWindowThreadProcessId(IntPtr h,out uint p);
 [DllImport("user32.dll")] public static extern bool GetGUIThreadInfo(uint t,ref GUI g);
 [DllImport("user32.dll")] public static extern bool PostMessage(IntPtr h,uint m,IntPtr w,IntPtr l);
 [DllImport("user32.dll")] public static extern bool EnumChildWindows(IntPtr h,EnumProc cb,IntPtr p);
 [DllImport("user32.dll")] public static extern int GetClassName(IntPtr h,StringBuilder s,int n);
 [DllImport("user32.dll")] public static extern int GetWindowText(IntPtr h,StringBuilder s,int n);
}
"@
$gameProc=Get-Process -Id $GameProcessId
if($gameProc.ProcessName -ne 'ACS 5') {throw 'Only the selected BR test process may receive input'}
$handle=$gameProc.MainWindowHandle
if($ButtonTitle) {
 $deadline=[DateTime]::UtcNow.AddSeconds(10)
 do {
 $buttons=[Collections.Generic.List[IntPtr]]::new()
 [ArcadeNative]::EnumChildWindows($handle,{param($h,$p) $class=[Text.StringBuilder]::new(200);$label=[Text.StringBuilder]::new(200); [ArcadeNative]::GetClassName($h,$class,200)|Out-Null;[ArcadeNative]::GetWindowText($h,$label,200)|Out-Null;if($class.ToString() -eq 'Button' -and $label.ToString() -eq $ButtonTitle){$buttons.Add($h)};return $true},[IntPtr]::Zero)|Out-Null
 if($buttons.Count -eq 0){Start-Sleep -Milliseconds 200}
 } while($buttons.Count -eq 0 -and [DateTime]::UtcNow -lt $deadline)
 if($buttons.Count -ne 1){throw "Expected one button '$ButtonTitle', found $($buttons.Count)"}
 [ArcadeNative]::PostMessage($buttons[0],0xF5,[IntPtr]::Zero,[IntPtr]::Zero)|Out-Null
} else {
 [uint32]$owner=0;$thread=[ArcadeNative]::GetWindowThreadProcessId($handle,[ref]$owner)
 $gui=New-Object ArcadeNative+GUI;$gui.cb=[Runtime.InteropServices.Marshal]::SizeOf($gui)
 [ArcadeNative]::GetGUIThreadInfo($thread,[ref]$gui)|Out-Null
 if($gui.focus -eq [IntPtr]::Zero -or $gameProc.MainWindowTitle -match '^(BR INVADERS|PADDLE PANIC|BYTE MUNCHER|CROSSWALK CHAOS|BITIPEDE|BITRIS)$') {
  $canvases=[Collections.Generic.List[IntPtr]]::new()
  [ArcadeNative]::EnumChildWindows($handle,{param($h,$p) $label=[Text.StringBuilder]::new(200);[ArcadeNative]::GetWindowText($h,$label,200)|Out-Null;if($label.ToString() -eq 'Winterm FullArea'){$canvases.Add($h)};return $true},[IntPtr]::Zero)|Out-Null
  if($canvases.Count -ne 1){throw 'No unique BR keyboard canvas'}
  $gui.focus=$canvases[0]
 }
 [uint32]$focusOwner=0;[ArcadeNative]::GetWindowThreadProcessId($gui.focus,[ref]$focusOwner)|Out-Null
 if($focusOwner -ne $GameProcessId){throw 'Focus control is outside test process'}
 foreach($character in $Text.ToCharArray()) {
   [ArcadeNative]::PostMessage($gui.focus,0x102,[IntPtr][int]$character,[IntPtr]1)|Out-Null
 }
}
Write-Output "Addressed native control in arcade PID $GameProcessId only."
