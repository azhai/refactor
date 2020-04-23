@ECHO OFF

REM USAGE:
REM     make      Rebuild the project
REM     make pre  Install all packages and rebuild the project

IF "%1" == "" GOTO all
IF "%1" == "all" GOTO all
IF "%1" == "build" GOTO build

:pre
    go.exe mod tidy
    go.exe mod vendor
:all
    del reverse.exe
:build
    go.exe build -ldflags="-s -w" -mod=vendor -o reverse.exe ./cmd/
