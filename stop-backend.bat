@echo off
echo Stopping PC Monitor Backend...
taskkill /F /IM pcmon.exe 2>nul
if %errorlevel% == 0 (
    echo Backend stopped successfully.
) else (
    echo No backend process found.
)
pause
