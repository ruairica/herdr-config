# One-time setup for a new machine: point herdr at this repo's config.toml, build pathpick,
# and make sure the tools the key bindings use are installed.
# Usage: pwsh -File install.ps1
$ErrorActionPreference = 'Stop'

$repo = $PSScriptRoot
$expected = 'C:\Source\personal\herdr-config'

# config.toml uses absolute paths (herdr has no relative ones), so the repo must live here.
if ($repo -ne $expected) {
    Write-Warning "Repo is at $repo but config.toml expects $expected. Move it, or edit the paths in config.toml."
}

# Point herdr at the repo config (user-level, read when the herdr server starts).
$config = Join-Path $repo 'config.toml'
[Environment]::SetEnvironmentVariable('HERDR_CONFIG_PATH', $config, 'User')
$env:HERDR_CONFIG_PATH = $config
Write-Host "HERDR_CONFIG_PATH = $config"

# Tools used by the key bindings.
$tools = @{
    fd      = 'sharkdp.fd'
    fzf     = 'junegunn.fzf'
    lazygit = 'JesseDuffield.lazygit'
    go      = 'GoLang.Go'
}
foreach ($name in $tools.Keys) {
    if (Get-Command $name -ErrorAction SilentlyContinue) {
        Write-Host "ok: $name"
    } else {
        Write-Host "installing $name ($($tools[$name]))"
        winget install --id $tools[$name] --exact --accept-source-agreements --accept-package-agreements
    }
}

# Build pathpick. A fresh winget install of Go isn't on this session's PATH yet.
$go = (Get-Command go -ErrorAction SilentlyContinue).Source
if (-not $go) { $go = 'C:\Program Files\Go\bin\go.exe' }
Push-Location (Join-Path $repo 'pathpick')
try {
    & $go build -ldflags="-s -w" -o pathpick.exe .
    if ($LASTEXITCODE) { throw "go build failed ($LASTEXITCODE)" }
    Write-Host 'built pathpick\pathpick.exe'
} finally {
    Pop-Location
}

Write-Host ''
Write-Host 'Done. Restart the herdr server so it picks up HERDR_CONFIG_PATH (reload-config is not enough).'
Write-Host 'New terminals may also be needed for fd/fzf/lazygit to be on PATH.'
