@echo off
title Hardware Monitor Backend Manager
color 0A

:menu
cls
echo ========================================
echo   Hardware Monitor Backend Manager
echo ========================================
echo.
echo 1. Check Status
echo 2. Start Backend (Hidden)
echo 3. Stop Backend
echo 4. View Logs (if running)
echo 5. Restart Backend
echo 6. Exit
echo.
echo ========================================
set /p choice="Enter your choice (1-6): "

if "%choice%"=="1" goto status
if "%choice%"=="2" goto start
if "%choice%"=="3" goto stop
if "%choice%"=="4" goto logs
if "%choice%"=="5" goto restart
if "%choice%"=="6" goto end
goto menu

:status
cls
echo Checking backend status...
echo.
tasklist | findstr "pcmon.exe" >nul
if %errorlevel% == 0 (
    echo [RUNNING] Backend is currently running
    echo.
    echo Process details:
    tasklist | findstr "pcmon.exe"
) else (
    echo [STOPPED] Backend is not running
)
echo.
pause
goto menu

:start
cls
echo Starting backend in hidden mode...
start "" "%~dp0start-backend-hidden.vbs"
timeout /t 2 >nul
echo.
echo Backend started! Use option 1 to verify.
echo.
pause
goto menu

:stop
cls
echo Stopping backend...
taskkill /F /IM pcmon.exe 2>nul
if %errorlevel% == 0 (
    echo [SUCCESS] Backend stopped
) else (
    echo [INFO] No backend process found
)
echo.
pause
goto menu

:logs
cls
echo Last 50 lines of backend log:
echo ========================================
echo.
if exist "%~dp0backend\pcmon.log" (
    powershell -Command "Get-Content '%~dp0backend\pcmon.log' -Tail 50"
) else (
    echo No log file found. Start the backend first.
)
echo.
pause
goto menu

:restart
cls
echo Stopping backend...
taskkill /F /IM pcmon.exe 2>nul
timeout /t 1 >nul
echo Starting backend in hidden mode...
start "" "%~dp0start-backend-hidden.vbs"
timeout /t 2 >nul
echo.
echo Done! Use option 1 to check status.
echo.
pause
goto menu

:end
exit
