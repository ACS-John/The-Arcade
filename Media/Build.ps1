param([string]$NodeExe='node',[string]$ModulesPath=$env:NODE_PATH)
$ErrorActionPreference='Stop'
Push-Location $PSScriptRoot
try {
    & go test ./...
    if($LASTEXITCODE -ne 0){throw 'Audio tests failed'}
    & go build -buildvcs=false -ldflags '-H windowsgui' -o ArcadeAudio.exe .
    if($LASTEXITCODE -ne 0){throw 'Audio build failed'}
    & ./ArcadeAudio.exe -render ../audio
    if($LASTEXITCODE -ne 0){throw 'WAV render failed'}
    if($ModulesPath){$env:NODE_PATH=$ModulesPath}
    & $NodeExe ./build-sprites.cjs
    if($LASTEXITCODE -ne 0){throw 'Sprite build failed; Node.js and sharp are required'}
} finally {Pop-Location}
