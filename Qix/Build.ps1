$ErrorActionPreference='Stop'
Push-Location $PSScriptRoot
try {
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw 'Qix tests failed' }
    go build -buildvcs=false -ldflags '-H windowsgui' -o Qix.exe .
    if ($LASTEXITCODE -ne 0) { throw 'Qix build failed. Close Qix before rebuilding.' }
} finally { Pop-Location }