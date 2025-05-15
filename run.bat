@echo off
echo Building WhatsApp Gateway...
go build
if %ERRORLEVEL% neq 0 (
    echo Build failed!
    pause
    exit /b %ERRORLEVEL%
)
echo Build successful! Starting server...
echo.
go-simple-whatsapp-gateway2.exe
pause