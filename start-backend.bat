@echo off
cd /d "%~dp0backend"

REM Check if built executable exists, otherwise build it
if not exist "pcmon.exe" (
    echo Building backend for first time...
    go build -o pcmon.exe main.go
)

REM Run the executable
pcmon.exe
