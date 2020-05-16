@ECHO OFF

REM USAGE:
REM     make      Rebuild the project

IF "%1" == "" GOTO all
IF "%1" == "all" GOTO all
IF "%1" == "build" GOTO build

:all
    del reverse.exe
:build
    go.exe build -ldflags="-s -w" -o reverse.exe ./cmd/
