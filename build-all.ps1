# Windows x64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o hlstatsz-forwarder.exe

# Linux x64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o hlstatsz-forwarder

# Reset
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host "✅ Build complete. Files: forwarder-windows.exe, forwarder-linux"
