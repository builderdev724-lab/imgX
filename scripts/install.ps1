$ErrorActionPreference = 'Stop'

$Repo = if ($env:IMGX_REPO) { $env:IMGX_REPO } else { 'priyojeet/imgx' }
$Version = if ($env:IMGX_VERSION) { $env:IMGX_VERSION } else { 'latest' }
$InstallDir = if ($env:IMGX_INSTALL_DIR) { $env:IMGX_INSTALL_DIR } else { Join-Path $HOME 'bin' }

$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
if ($Version -eq 'latest') {
  $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
  $Tag = $release.tag_name
} else {
  $Tag = $Version
}

$Zip = "imgx_$($Tag.TrimStart('v'))_windows_$Arch.zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$Zip"
$Tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("imgx-" + [guid]::NewGuid())
New-Item -ItemType Directory -Path $Tmp | Out-Null
try {
  $ZipPath = Join-Path $Tmp $Zip
  $Checksums = "imgx_$($Tag.TrimStart('v'))_checksums.txt"
  $ChecksumsPath = Join-Path $Tmp $Checksums
  Invoke-WebRequest -Uri $Url -OutFile $ZipPath
  Invoke-WebRequest -Uri "https://github.com/$Repo/releases/download/$Tag/$Checksums" -OutFile $ChecksumsPath
  $line = Get-Content $ChecksumsPath | Where-Object { $_ -match "\s$([regex]::Escape($Zip))$" } | Select-Object -First 1
  if (-not $line) { throw "Checksum entry not found for $Zip" }
  $expected = ($line -split '\s+')[0]
  $actual = (Get-FileHash -Algorithm SHA256 -Path $ZipPath).Hash.ToLowerInvariant()
  if ($expected.ToLowerInvariant() -ne $actual) { throw "SHA-256 verification failed" }
  Expand-Archive -Path $ZipPath -DestinationPath $Tmp -Force
  New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
  Copy-Item (Join-Path $Tmp 'imgx.exe') (Join-Path $InstallDir 'imgx.exe') -Force
  Write-Host "imgx installed to $(Join-Path $InstallDir 'imgx.exe')"
  & (Join-Path $InstallDir 'imgx.exe') --version
  Write-Host "Add $InstallDir to PATH if it is not already there."
} finally {
  Remove-Item -Recurse -Force $Tmp -ErrorAction SilentlyContinue
}
