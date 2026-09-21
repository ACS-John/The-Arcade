param([Parameter(Mandatory=$true)][int]$GameProcessId,[Parameter(Mandatory=$true)][string]$OutFile)
Add-Type -AssemblyName System.Drawing
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class ArcadeCapture {
 [StructLayout(LayoutKind.Sequential)] public struct RECT { public int Left,Top,Right,Bottom; }
 [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out RECT r);
 [DllImport("user32.dll")] public static extern bool PrintWindow(IntPtr h, IntPtr dc, uint flags);
}
"@
$gameProc=Get-Process -Id $GameProcessId
if($gameProc.ProcessName -notin @('ACS 5','Qix')) { throw 'Only capture the selected BR test process' }
$rect=New-Object ArcadeCapture+RECT
[ArcadeCapture]::GetWindowRect($gameProc.MainWindowHandle,[ref]$rect) | Out-Null
$bmp=New-Object Drawing.Bitmap ($rect.Right-$rect.Left),($rect.Bottom-$rect.Top)
$graphics=[Drawing.Graphics]::FromImage($bmp)
$dc=$graphics.GetHdc()
try { $captured=[ArcadeCapture]::PrintWindow($gameProc.MainWindowHandle,$dc,2) } finally {$graphics.ReleaseHdc($dc)}
if(-not $captured) {throw 'PrintWindow failed'}
$bmp.Save($OutFile,[Drawing.Imaging.ImageFormat]::Png)
$graphics.Dispose(); $bmp.Dispose()
Write-Output $OutFile
