@echo off
setlocal EnableDelayedExpansion

echo.
echo ====================================
echo  📝 Create New Migration
echo ====================================
echo.

set /p NAME="Nama migrasi (contoh: create_users_table): "

if "%NAME%"=="" (
    echo ❌ Nama migrasi tidak boleh kosong!
    pause
    exit /b 1
)

echo.
echo 📝 Membuat file migrasi: %NAME%
echo.

docker run --rm ^
-v "%cd%/db/migration:/migrations" ^
migrate/migrate:v4.19.1 ^
create -ext sql -dir /migrations -seq %NAME%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ File migrasi berhasil dibuat di db/migration/
) else (
    echo.
    echo ❌ Gagal membuat file migrasi!
)

echo.
pause