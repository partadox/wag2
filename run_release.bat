@echo off
echo Building WhatsApp Gateway in release mode...
set GIN_MODE=release
go build -ldflags="-s -w" -o go-simple-whatsapp-gateway2-release.exe
if %ERRORLEVEL% neq 0 (
    echo Build failed!
    pause
    exit /b %ERRORLEVEL%
)
echo Build successful! Starting server in release mode...
echo.
echo =========================================
echo Go Simple WhatsApp Gateway (Release Mode)
echo =========================================
echo.
echo Access the web UI at: http://localhost:8080
echo Default API Key: changeme (from .env file)
echo.
go-simple-whatsapp-gateway2-release.exe
pause