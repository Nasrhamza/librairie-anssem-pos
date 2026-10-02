@echo off
setlocal
cd /d "%~dp0"
where go >nul 2>nul
if errorlevel 1 (
  echo [ERREUR] Go n'est pas installe ou n'est pas dans PATH.
  pause
  exit /b 1
)
go run .
