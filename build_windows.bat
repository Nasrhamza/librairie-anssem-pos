@echo off
setlocal
cd /d "%~dp0"
set "GO_EXE=go"
where go >nul 2>nul
if errorlevel 1 if exist "%LOCALAPPDATA%\CodexToolchains\go1.27.1\go\bin\go.exe" set "GO_EXE=%LOCALAPPDATA%\CodexToolchains\go1.27.1\go\bin\go.exe"
"%GO_EXE%" version >nul 2>nul
if errorlevel 1 (
  echo [ERREUR] Go n'est pas installe ou n'est pas dans PATH.
  echo Installer Go depuis https://go.dev/dl/ puis relancer ce fichier.
  pause
  exit /b 1
)
if not exist dist mkdir dist
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
"%GO_EXE%" test ./...
if errorlevel 1 (
  echo [ERREUR] Les tests ont echoue. Build annule.
  pause
  exit /b 1
)
"%GO_EXE%" build -trimpath -ldflags="-s -w -H windowsgui" -o dist\Jamel-v1.exe .
if errorlevel 1 (
  echo [ERREUR] Build Windows echoue.
  pause
  exit /b 1
)
if not exist installer\payload mkdir installer\payload
copy /Y dist\Jamel-v1.exe installer\payload\Jamel-v1.exe >nul
"%GO_EXE%" build -trimpath -ldflags="-s -w -H windowsgui" -o dist\Jamel-v1-Setup.exe .\installer
if errorlevel 1 (
  echo [ERREUR] Build du Setup Windows echoue.
  pause
  exit /b 1
)
"%GO_EXE%" build -trimpath -ldflags="-s -w -H windowsgui" -o owner-tools\Jamel-License-Generator.exe .\owner-tools\license-generator
if errorlevel 1 (
  echo [ERREUR] Build du generateur prive de licences echoue.
  pause
  exit /b 1
)
echo.
echo [OK] EXE cree : dist\Jamel-v1.exe
echo [OK] Setup cree : dist\Jamel-v1-Setup.exe
echo [PRIVE] Generateur cree : owner-tools\Jamel-License-Generator.exe
echo [ATTENTION] Ne jamais envoyer le dossier owner-tools au client.
pause
