# pathpick

A [herdr](https://herdr.dev/tui/) popup that fuzzy-picks a folder and types its full path into the pane you were in, without pressing Enter. It's handy for dropping a repo path into an agent prompt mid-sentence.

## Requirements

- [herdr](https://herdr.dev/tui/) 0.9.5+
- [fd](https://github.com/sharkdp/fd) and [fzf](https://github.com/junegunn/fzf) on your `PATH`
- Go, to build

## Build

```
go build -ldflags="-s -w" -o pathpick.exe .
```

## Set up

Add a popup binding to herdr's `config.toml`, then run `herdr server reload-config`:

```toml
[[keys.command]]
key = "prefix+f"
type = "popup"
command = 'C:\path\to\pathpick.exe -root C:\Source'
width = "80%"
height = "60%"
```

## Usage

Press the key, pick a folder, and press Enter. The path is typed into the pane with a space on each side, and wrapped in quotes if it contains spaces. Press Ctrl+O instead to open the folder in a new herdr workspace and switch to it. Esc cancels.

`-root` sets the folder whose subfolders are listed (one level deep). Without it, pathpick lists the subfolders of the focused pane's folder.
