@echo off
echo Building hook DLL...

REM 使用MinGW-w64编译DLL
echo Compiling with MinGW-w64...
gcc -shared -o hook_dll.dll hook_dll.c -lkernel32 -luser32 -lcomctl32

if %ERRORLEVEL% EQU 0 (
    echo DLL built successfully: hook_dll.dll
    echo.
    echo Testing DLL exports...
    dumpbin /exports hook_dll.dll 2>nul || echo dumpbin not available, skipping export test
) else (
    echo Failed to build DLL with MinGW-w64
    echo.
    echo Trying Visual Studio compiler...
    cl /LD hook_dll.c /Fe:hook_dll.dll /link kernel32.lib user32.lib comctl32.lib
    
    if %ERRORLEVEL% EQU 0 (
        echo DLL built successfully with Visual Studio: hook_dll.dll
    ) else (
        echo Failed to build DLL with both compilers
        echo Please install MinGW-w64 or Visual Studio Build Tools
    )
)

echo.
echo DLL compilation completed.
pause
