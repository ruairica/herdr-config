# pathpick

Herdr popup that fuzzy-picks a folder directly under `-root` (default: the focused pane's folder) and types its path into the focused pane, or on Ctrl+O opens it in a new focused workspace (`herdr workspace create --focus`). Bound to `prefix+f` in the repo-root `config.toml` (`type = "popup"`) with `-root C:\Source`.

## Build

```
go build -ldflags="-s -w" -o pathpick.exe .
```

## Debug log

Each run appends start/fd-done/fzf-done timings to `%LOCALAPPDATA%\pathpick.log` (reset past 1 MB), added to diagnose popups that come up blank until clicked.
