# go-tdlib

Go wrapper for [TDLib (Telegram Database Library)](https://github.com/tdlib/td) with full support of TDLib v1.8.0

## TDLib installation

Use [TDLib build instructions](https://tdlib.github.io/td/build.html) with checkmarked `Install built TDLib to /usr/local instead of placing the files to td/tdlib`.

### Windows

Build with environment variables:

在写全路径，不能是相对路径
```
CGO_CFLAGS=-IC:/path/to/tdlib/build/tdlib/include
CGO_LDFLAGS=-LC:/path/to/tdlib/build/tdlib/bin -ltdjson
```

Example for PowerShell:

```powershell
$env:CGO_CFLAGS="-IC:/td/tdlib/include"; $env:CGO_LDFLAGS="-LC:/td/tdlib/bin -ltdjson"; go build -trimpath -ldflags="-s -w" -o demo.exe .\cmd\demo.go

$env:CGO_CFLAGS="-IC:/path/to/tdlib/include"; $env:CGO_LDFLAGS="-LC:/path/to/tdlib/bin -ltdjson"; go build -trimpath -ldflags="-s -w" -o demo.exe .\cmd\demo.go

```