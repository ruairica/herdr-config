# herdr-config

My [herdr](https://herdr.dev/tui/) config, helper scripts and the pathpick popup, kept in one place.

Herdr reads `config.toml` straight from this repo via the `HERDR_CONFIG_PATH` environment variable, so there's nothing to sync: edit, reload, commit.

## Layout

```
config.toml               the live herdr config (HERDR_CONFIG_PATH points here)
scripts/new-agent-tab.ps1 new tab + launch an agent CLI (prefix+alt+c / prefix+alt+p)
pathpick/                 fuzzy folder picker popup (prefix+f), see pathpick/README.md
install.ps1               one-time setup for a new machine
```

## Set up a machine

Clone to `C:\Source\personal\herdr-config` (the config uses absolute paths), then:

```powershell
pwsh -File install.ps1
```

It sets `HERDR_CONFIG_PATH` at user level, installs `fd`, `fzf`, `lazygit` and Go with `winget` if missing, and builds `pathpick\pathpick.exe` (gitignored). Restart the herdr server afterwards; it only reads `HERDR_CONFIG_PATH` at startup, so `herdr server reload-config` isn't enough the first time.

## Day to day

Edit `config.toml`, run `herdr server reload-config`, commit when happy.

After changing pathpick, rebuild it:

```powershell
cd pathpick; go build -ldflags="-s -w" -o pathpick.exe .
```

## Notes

- Only my own files live here. `%APPDATA%\herdr` still holds herdr's runtime state (sockets, logs, `session.json`, snapshots), which stays out of git.
- `herdr config reset-keys` rewrites the config file, which is this repo's file. Git shows the diff, so it's easy to revert.
- Fallback if the env var ever misbehaves: symlink `%APPDATA%\herdr\config.toml` to this repo's file (needs Developer Mode or an admin shell on Windows).
