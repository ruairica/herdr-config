# Create a new Herdr tab and immediately launch a command (e.g. an agent CLI) in it.
# Usage: new-agent-tab.ps1 <command> [tab-label]
param(
    [string]$Command = 'claude',
    [string]$Label = ''
)

$ErrorActionPreference = 'Stop'
if (-not $Label) { $Label = ($Command -split '\s+')[0] }

# Inherit the cwd of whichever pane is focused when the keybinding fires.
$cwd = $null
try { $cwd = (& herdr pane current | ConvertFrom-Json).result.pane.cwd } catch { }

$args = @('tab', 'create', '--label', $Label, '--focus')
if ($cwd) { $args += @('--cwd', $cwd) }

$tab = & herdr @args | ConvertFrom-Json
$pane = $tab.result.root_pane.pane_id
& herdr pane run $pane $Command | Out-Null
