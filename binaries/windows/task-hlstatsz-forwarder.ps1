$exePath = Join-Path $PSScriptRoot "hlstatsz-forwarder.exe"
while ($true) {
    Start-Process -FilePath $exePath -Wait
    Start-Sleep -Seconds 1  # brief pause before relaunch
}
