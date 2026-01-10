@echo off
setlocal EnableDelayedExpansion

echo.
echo ====================================
echo  🔄 Database Rollback Tool
echo ====================================
echo.

set /p STEPS="Berapa step yang ingin di-rollback? (default: 1): "
if "%STEPS%"=="" set STEPS=1

echo.
echo 🔄 Rollback %STEPS% migrasi terakhir...
echo.

docker run --rm ^
-v "%cd%/db/migration:/migrations" ^
--network pos-api_default ^
migrate/migrate:v4.19.1 ^
-path=/migrations ^
-database "postgres://root:secret@pos_postgres:5432/pos_db?sslmode=disable" ^
-verbose down %STEPS%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Rollback berhasil!
) else (
    echo.
    echo ❌ Rollback gagal!
)

echo.
pause