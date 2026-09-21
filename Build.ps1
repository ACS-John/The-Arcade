param([switch]$Test,[string]$Compiler=$env:ARCADE_COMPILER,[string]$BrExe=$env:ARCADE_BR_EXE,[string]$BrlsExe=$env:ARCADE_BRLS)
$ErrorActionPreference='Stop'
$projectRoot=$PSScriptRoot
$acsRoot=Join-Path (Split-Path -Parent $projectRoot) 'Dev-5'
if(-not $Compiler){$Compiler=Join-Path $acsRoot 'Nine Tailed Fox/nineTailedFox.exe'}
if(-not $BrExe){$BrExe=Join-Path $acsRoot 'ACS 5.exe'}
if(-not $BrlsExe){$BrlsExe=Join-Path $acsRoot 'context/dev/tools/brls.exe'}
if(-not(Test-Path -LiteralPath $Compiler)){throw 'Set ARCADE_COMPILER or pass -Compiler with your Nine-Tailed Fox compiler path.'}
if($Test -and -not(Test-Path -LiteralPath $BrExe)){throw 'Set ARCADE_BR_EXE or pass -BrExe with your licensed BR runtime path.'}
Push-Location $projectRoot
try {
 New-Item -ItemType Directory -Path .build -Force | Out-Null
 $sources=@('Arcade.panda','BR Invaders/BR Invaders.panda','Paddle Panic/Paddle Panic.panda','Byte Muncher/Byte Muncher.panda','Crosswalk Chaos/Crosswalk Chaos.panda','Bitipede/Bitipede.panda','Bitris/Bitris.panda','Aphelion/Aphelion.panda')
 $testNames=@('Invaders','Paddle','Muncher','Crosswalk','Bitipede')
 if($Test){$sources+=@($testNames|ForEach-Object{'Tests/'+$_+'.panda'});$sources+='Bitris/Tests.panda'}
 $sources|ForEach-Object{Join-Path $projectRoot $_}|Set-Content .build/build-files.txt -Encoding ascii
 & $Compiler ('filelist:'+(Join-Path $projectRoot '.build/build-files.txt')) /ai /leaveBrs
 if($LASTEXITCODE -ne 0){throw 'Compilation failed'}
 if(Test-Path -LiteralPath $BrlsExe){
  foreach($source in $sources){
   & $BrlsExe -check -sema ([IO.Path]::ChangeExtension($source,'.brs'))
   if($LASTEXITCODE -ne 0){throw "Syntax/semantic check failed: $source"}
  }
 } else {Write-Warning 'BRLS not configured; compiler validation completed, optional static check skipped.'}
 if($Test){
  @('LOGGING 10, .build\br.log, unattended','Application_Name Arcade Tests','WorkPath %Temp%','Drive S,.,X,\','Console Off','WorkStack 99999999','RPNStack 2000','FileNames Mixed_Case Search','Font Courier New','Option 26 On','Option 46 On','Date never') | Set-Content .build/headless.sys -Encoding ascii
  function Invoke-ArcadeProc([string[]]$Statements){
   $Statements|Set-Content .build/run.prc -Encoding ascii
   if(Test-Path .build/br.log){Remove-Item -LiteralPath .build/br.log}
   & $BrExe 'PROC .build\run.prc' '-.build\headless.sys' | Out-Null
   if(-not(Test-Path .build/br.log)){throw 'BR did not produce a runtime log'}
   $log=Get-Content .build/br.log -Raw
   if($log -match 'Unattended processing terminated by error'){throw $log}
  }
  foreach($name in @($testNames)+@('Bitris')){
   $program=if($name -eq 'Bitris'){'S:\Bitris\Tests.br'}else{'S:\Tests\'+$name+'.br'}
   $report=if($name -eq 'Bitris'){'Bitris/test-results.txt'}else{'Tests/'+$name+'-results.txt'}
   if(Test-Path -LiteralPath $report){Remove-Item -LiteralPath $report}
   Invoke-ArcadeProc @(('load "'+$program+'"'),'run','execute "system"')
   $result=Get-Content -LiteralPath $report -Raw
   if($result -match 'FAIL:' -or $result -notmatch 'RESULT: \d+ checks, 0 failures'){throw "Failed: $name $result"}
   Write-Output ($name+': '+$result.Trim())
  }
  foreach($game in @('BR Invaders','Paddle Panic','Byte Muncher','Crosswalk Chaos','Bitipede','Bitris')){
   Invoke-ArcadeProc @('execute ''let setenv("ArcadeSmoke","Yes")''',('load "S:\'+$game+'\'+$game+'.br"'),'run','execute "system"')
   Write-Output "Rendered startup passed: $game"
  }
 }
 Write-Output 'Arcade build complete.'
} finally {Pop-Location}
