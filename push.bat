@echo off

REM Cek apakah commit message diisi
if "%~1"=="" (
    echo Commit message kosong!
    echo Contoh: push.bat "update login feature"
    exit /b 1
)

echo Git add...
git add .

echo Git commit...
git commit -m "%~1"

echo Git push...
git push origin main

echo Push selesai!
