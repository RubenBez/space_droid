CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o build/SpaceDroid-win.exe -ldflags -s -w
