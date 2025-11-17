@echo off
echo Building backend...
cd /d "%~dp0backend"
go build -o pcmon.exe main.go
if %errorlevel% == 0 (
    echo Build successful! Executable created: pcmon.exe
) else (
    echo Build failed!
)
pause
