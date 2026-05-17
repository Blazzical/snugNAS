; snugNAS Windows installer
;
; Build with NSIS 3.x (https://nsis.sourceforge.io/):
;   makensis installer\snugnas.nsi
;
; Expects a freshly built snugnas.exe at the repo root. See docs\build.md.

!define APP_NAME      "snugNAS"
!ifndef APP_VERSION
  !define APP_VERSION "0.0.1-dev"
!endif
!define APP_PUBLISHER "snugNAS contributors"
!define APP_URL       "https://github.com/Blazzical/snugNAS"
!define APP_EXE       "snugnas.exe"
!define APP_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"

Name "${APP_NAME}"
OutFile "snugnas-setup.exe"
InstallDir "$PROGRAMFILES64\${APP_NAME}"
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallDir"
RequestExecutionLevel admin
SetCompressor /SOLID lzma
Unicode true

; VIProductVersion needs an x.y.z.b 4-part numeric form. The CLI passes a
; semver like "0.8.1"; map any missing trailing parts to .0 so e.g. 0.8.1
; -> 0.8.1.0. Built-in cleanups would be nicer in NSIS 4 but for now we
; just require the caller to pass an x.y.z value.
!ifndef APP_VERSION_NUMERIC
  !define APP_VERSION_NUMERIC "${APP_VERSION}.0"
!endif

VIProductVersion "${APP_VERSION_NUMERIC}"
VIAddVersionKey "ProductName"     "${APP_NAME}"
VIAddVersionKey "FileDescription" "snugNAS installer"
VIAddVersionKey "FileVersion"     "${APP_VERSION}"
VIAddVersionKey "ProductVersion"  "${APP_VERSION}"
VIAddVersionKey "CompanyName"     "${APP_PUBLISHER}"
VIAddVersionKey "LegalCopyright"  "MIT"

!include "MUI2.nsh"

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES

!define MUI_FINISHPAGE_RUN "$INSTDIR\${APP_EXE}"
!define MUI_FINISHPAGE_RUN_PARAMETERS "tray"
!define MUI_FINISHPAGE_RUN_TEXT "Launch snugNAS now"

; Optional: hint about autostart in the finish page. Keep it advisory so
; we don't write to HKCU\Run from an admin installer (different user scope).
!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Run `snugnas autostart enable` to start at login"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION ShowAutostartHint
!define MUI_FINISHPAGE_SHOWREADME_NOTCHECKED
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Function ShowAutostartHint
  MessageBox MB_OK \
    "After install, run `snugnas autostart enable` from an ordinary PowerShell$\n\
to have snugNAS launch on login. Don't run it from this elevated installer$\n\
— Windows would register autostart for the SYSTEM account, not you."
FunctionEnd

Section "snugNAS" SecCore
  SectionIn RO  ; required component

  SetOutPath "$INSTDIR"
  File "..\${APP_EXE}"
  File "..\LICENSE"
  File "..\README.md"

  ; Start-menu shortcut
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortCut  "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"        "$INSTDIR\${APP_EXE}" "tray"   "$INSTDIR\${APP_EXE}" 0
  CreateShortCut  "$SMPROGRAMS\${APP_NAME}\Wizard.lnk"             "$INSTDIR\${APP_EXE}" "wizard" "$INSTDIR\${APP_EXE}" 0
  CreateShortCut  "$SMPROGRAMS\${APP_NAME}\Uninstall ${APP_NAME}.lnk" "$INSTDIR\uninstall.exe"

  ; Track install location for upgrades
  WriteRegStr HKLM "Software\${APP_NAME}" "InstallDir" "$INSTDIR"

  ; Add/Remove Programs entry
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayName"     "${APP_NAME}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayVersion"  "${APP_VERSION}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "Publisher"       "${APP_PUBLISHER}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "URLInfoAbout"    "${APP_URL}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "DisplayIcon"     "$INSTDIR\${APP_EXE}"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "${APP_UNINST_KEY}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
  WriteRegDWORD HKLM "${APP_UNINST_KEY}" "NoModify" 1
  WriteRegDWORD HKLM "${APP_UNINST_KEY}" "NoRepair" 1

  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\LICENSE"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\uninstall.exe"
  RMDir  "$INSTDIR"

  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\Wizard.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\Uninstall ${APP_NAME}.lnk"
  RMDir  "$SMPROGRAMS\${APP_NAME}"

  DeleteRegKey HKLM "${APP_UNINST_KEY}"
  DeleteRegKey HKLM "Software\${APP_NAME}"

  ; Note: we deliberately do NOT delete the user's config or storage
  ; (%APPDATA%\snugNAS, plus whatever they set as storage_root). The
  ; uninstall removes the binary only — data preservation is the safer
  ; default.
SectionEnd
