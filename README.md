Metabot
========

Vanilla WoW bot framework, exclusively for 1.12.1 build 5875.

## Building

Requirements:
* Windows 10/11
* go 1.22
* TDM GCC (32bit)
* protoc 26.1

Protobuf:
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/game.proto
```

Injector:
```
cd cmd/injector
go build .
```

dll:
```
cd cmd/dll
go build -o metabot.dll -buildmode=c-shared
```

## Running

The project currently requires WoW to be installed to `C:\wow` and the injector looks for `metabot.dll` at 
`C:\metabot\cmd\dll\metabot.dll`. You can edit this in `cmd/injector/main.go`.

After setting things up, you should be able to simply run `injector.exe` as an administrator and things _should_ work.
