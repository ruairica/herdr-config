// pathpick: fuzzy-pick a folder directly under -root and type its path into the herdr pane that
// was focused when the popup opened (no Enter), or with Ctrl+O open it in a new herdr workspace.
// Run from a herdr key binding with type = "popup",
// which sets HERDR_ACTIVE_PANE_ID. Uses the same fd flags as Select-Folder in the pwsh profile.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	rootFlag := flag.String("root", ".", "folder whose subfolders are listed; a herdr popup starts in the focused pane's folder")
	flag.Parse()

	root, err := filepath.Abs(*rootFlag)
	if err != nil {
		fail(err.Error())
	}

	target := os.Getenv("HERDR_ACTIVE_PANE_ID")
	if target == "" {
		fail("no HERDR_ACTIVE_PANE_ID; run pathpick from a herdr popup key binding")
	}

	key, choice, ok := pick(root)
	if !ok {
		return // Esc, Ctrl+C or no match
	}

	path := filepath.Join(root, choice) // also drops fd's trailing backslash

	herdr := os.Getenv("HERDR_BIN_PATH")
	if herdr == "" {
		herdr = "herdr"
	}

	if key == "ctrl-o" {
		// Open the folder in a new workspace and switch to it; the popup closes when we exit.
		if out, err := exec.Command(herdr, "workspace", "create", "--cwd", path, "--label", filepath.Base(path), "--focus").CombinedOutput(); err != nil {
			fail(fmt.Sprintf("herdr workspace create failed: %v\n%s", err, out))
		}
		return
	}

	if strings.Contains(path, " ") {
		path = `"` + path + `"`
	}
	if out, err := exec.Command(herdr, "pane", "send-text", target, " "+path+" ").CombinedOutput(); err != nil {
		fail(fmt.Sprintf("herdr pane send-text failed: %v\n%s", err, out))
	}
}

// pick pipes fd's list of folders under root into fzf and returns the key that accepted the
// selection ("" for Enter, "ctrl-o") and the selected line.
func pick(root string) (string, string, bool) {
	fd := exec.Command("fd", "--type", "d", "--max-depth", "1", "--hidden", "--no-ignore", "--exclude", ".git", "--base-directory", root)
	// load:clear-screen redraws once fd's list is in; the popup sometimes comes up blank until clicked.
	fzf := exec.Command("fzf", "--header", "enter: insert path | ctrl-o: open in new workspace", "--expect", "ctrl-o", "--bind", "load:clear-screen")

	r, w, err := os.Pipe()
	if err != nil {
		fail(err.Error())
	}
	var out bytes.Buffer
	fd.Stdout = w
	fzf.Stdin, fzf.Stdout, fzf.Stderr = r, &out, os.Stderr

	start := time.Now()
	since := func() string { return time.Since(start).Round(time.Millisecond).String() }

	if err := fd.Start(); err != nil {
		fail("starting fd: " + err.Error())
	}
	if err := fzf.Start(); err != nil {
		fail("starting fzf: " + err.Error())
	}
	logf("start pane=%s root=%s fd_pid=%d fzf_pid=%d after=%s", os.Getenv("HERDR_ACTIVE_PANE_ID"), root, fd.Process.Pid, fzf.Process.Pid, since())
	// Drop our copies so fd sees a closed pipe if fzf exits early.
	r.Close()
	w.Close()

	// fd normally finishes long before fzf, so wait for it separately to log when the list was complete.
	fdDone := make(chan struct{})
	go func() {
		fdErr := fd.Wait()
		logf("fd done err=%v after=%s", fdErr, since())
		close(fdDone)
	}()

	fzfErr := fzf.Wait()
	<-fdDone

	// With --expect, fzf prints the accepting key on the first line (empty for Enter), then the selection.
	key, choice, _ := strings.Cut(strings.TrimRight(out.String(), "\r\n"), "\n")
	key, choice = strings.TrimSpace(key), strings.TrimSpace(choice)
	logf("fzf done err=%v key=%q choice=%q after=%s", fzfErr, key, choice, since())
	return key, choice, fzfErr == nil && choice != ""
}

// logf appends a timestamped line to pathpick.log in %LOCALAPPDATA%, to diagnose popups that come up
// blank. The log starts over once it passes 1 MB. Errors are ignored; logging must never break the picker.
func logf(format string, args ...any) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return
	}
	path := filepath.Join(dir, "pathpick.log")
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if fi, err := os.Stat(path); err == nil && fi.Size() > 1<<20 {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s pid=%d %s\n", time.Now().Format("2006-01-02 15:04:05.000"), os.Getpid(), fmt.Sprintf(format, args...))
}

// fail shows the message long enough to read before the popup closes.
func fail(msg string) {
	fmt.Fprintln(os.Stderr, "pathpick:", msg)
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
