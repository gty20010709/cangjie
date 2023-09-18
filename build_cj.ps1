Write-Output "Start Building..."

Write-Output "Building for Linux amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "linux"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/cj_linux_amd64/cj
tar -zcvf dist/cj_linux_amd64.tar.gz dist/cj_linux_amd64
Write-Output "Building for darwin amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "darwin"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/cj_darwin_amd64/cj
tar -zcvf dist/cj_darwin_amd64.tar.gz dist/cj_darwin_amd64
Write-Output "Building for windows amd64..."
$ENV:CGO_ENABLED = 0; $ENV:GOOS = "windows"; $ENV:GOARCH = "amd64"; go build --ldflags "-s -w" -o dist/cj_windows_amd64/cj.exe

$compress = @{
    Path             = "dist/cj_windows_amd64"
    CompressionLevel = "Fastest"
    DestinationPath  = "dist/cj_windows_amd64.zip"
    Force            = $true
}
Compress-Archive @compress


Write-Output "Build and Compress Complete. Cleaning up..."
# force remove
Remove-Item "dist/cj_linux_amd64" -Force -Recurse
Write-Output "rm dist/cj_linux_amd64" 
Remove-Item dist/cj_darwin_amd64 -Force -Recurse
Write-Output "rm dist/cj_darwin_amd64"
Remove-Item dist/cj_windows_amd64 -Force -Recurse
Write-Output "rm dist/cj_windows_amd64"
