@echo off
echo Building WhatsApp Gateway with debug logging...
go build -o go-simple-whatsapp-gateway2-debug.exe
if %ERRORLEVEL% neq 0 (
    echo Build failed!
    pause
    exit /b %ERRORLEVEL%
)
echo Build successful! Starting server with debug logging...
echo.
set GIN_MODE=debug
go-simple-whatsapp-gateway2-debug.exe
pause