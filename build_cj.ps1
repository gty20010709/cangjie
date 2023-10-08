$AppName = "cj"


Write-Output "Start Building..."

Write-Output "Building for Linux amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "linux"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/$($AppName)_linux_amd64/$($AppName)
tar -zcvf dist/$($AppName)_linux_amd64.tar.gz dist/$($AppName)_linux_amd64
Write-Output "Building for darwin amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "darwin"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/$($AppName)_darwin_amd64/$($AppName)
tar -zcvf dist/$($AppName)_darwin_amd64.tar.gz dist/$($AppName)_darwin_amd64
Write-Output "Building for windows amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "windows"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/$($AppName)_windows_amd64/$($AppName).exe

$compress = @{
    Path             = "dist/$($AppName)_windows_amd64"
    CompressionLevel = "Fastest"
    DestinationPath  = "dist/$($AppName)_windows_amd64.zip"
    Force            = $true
}
Compress-Archive @compress


Write-Output "Build and Compress Complete. Cleaning up..."
# force remove
Remove-Item "dist/$($AppName)_linux_amd64" -Force -Recurse
Write-Output "rm dist/$($AppName)_linux_amd64" 
Remove-Item dist/$($AppName)_darwin_amd64 -Force -Recurse
Write-Output "rm dist/$($AppName)_darwin_amd64"
Remove-Item dist/$($AppName)_windows_amd64 -Force -Recurse
Write-Output "rm dist/$($AppName)_windows_amd64"
