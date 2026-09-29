#ifndef AppVersion
  #define AppVersion "0.2.3"
#endif

[Setup]
AppId={{C5065AB5-3483-4B54-AF9E-1710338768F4}
AppName=AgentBoard
AppVersion={#AppVersion}
AppPublisher=AgentBoard
DefaultDirName=C:\AI\AgentBoard
DefaultGroupName=AgentBoard
OutputDir=..\release
OutputBaseFilename=agentboard-{#AppVersion}-windows-amd64-setup
SetupIconFile=..\web\public\agentboard.ico
LicenseFile=..\LICENSE
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
ChangesEnvironment=yes
MinVersion=10.0
UninstallDisplayIcon={app}\agentboard.ico

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "russian"; MessagesFile: "compiler:Languages\Russian.isl"

[Files]
Source: "..\release\stage\*"; DestDir: "{app}"; Flags: recursesubdirs createallsubdirs ignoreversion

[Icons]
Name: "{autoprograms}\AgentBoard"; Filename: "{app}\agentboard.exe"; Parameters: "open"; IconFilename: "{app}\agentboard.ico"
Name: "{autodesktop}\AgentBoard"; Filename: "{app}\agentboard.exe"; Parameters: "open"; IconFilename: "{app}\agentboard.ico"; Tasks: desktopicon

[Tasks]
Name: desktopicon; Description: "Create a desktop shortcut"; GroupDescription: "Additional shortcuts:"

[Code]
const
  EnvironmentKey = 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment';

function PathContains(const CurrentPath, Directory: String): Boolean;
begin
  Result := Pos(';' + Uppercase(Directory) + ';', ';' + Uppercase(CurrentPath) + ';') > 0;
end;

procedure AddToPath;
var
  CurrentPath, Directory: String;
begin
  Directory := ExpandConstant('{app}');
  if RegQueryStringValue(HKLM, EnvironmentKey, 'Path', CurrentPath) then begin
    if not PathContains(CurrentPath, Directory) then begin
      if (Length(CurrentPath) > 0) and (CurrentPath[Length(CurrentPath)] <> ';') then
        CurrentPath := CurrentPath + ';';
      if not RegWriteExpandStringValue(HKLM, EnvironmentKey, 'Path', CurrentPath + Directory) then
        RaiseException('Could not add AgentBoard to PATH.');
    end;
  end else if not RegWriteExpandStringValue(HKLM, EnvironmentKey, 'Path', Directory) then
    RaiseException('Could not create PATH.');
end;

procedure RemoveFromPath;
var
  CurrentPath, Entry, Directory, ResultPath: String;
  Separator: Integer;
begin
  if not RegQueryStringValue(HKLM, EnvironmentKey, 'Path', CurrentPath) then Exit;
  Directory := ExpandConstant('{app}');
  ResultPath := '';
  while Length(CurrentPath) > 0 do begin
    Separator := Pos(';', CurrentPath);
    if Separator = 0 then begin
      Entry := CurrentPath;
      CurrentPath := '';
    end else begin
      Entry := Copy(CurrentPath, 1, Separator - 1);
      Delete(CurrentPath, 1, Separator);
    end;
    if (Entry <> '') and (CompareText(Entry, Directory) <> 0) then begin
      if ResultPath <> '' then ResultPath := ResultPath + ';';
      ResultPath := ResultPath + Entry;
    end;
  end;
  RegWriteExpandStringValue(HKLM, EnvironmentKey, 'Path', ResultPath);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then AddToPath;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then RemoveFromPath;
end;
