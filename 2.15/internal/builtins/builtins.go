package builtins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func IsBuiltin(name string) bool {
	switch name {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	default:
		return false
	}
}

func Run(args []string, in io.Reader, out io.Writer) error {
	if len(args) == 0 {
		return nil
	}
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return fmt.Errorf("cd: missing path")
		}
		return os.Chdir(args[1])
	case "pwd":
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, cwd)
		return err
	case "echo":
		_, err := fmt.Fprintln(out, strings.Join(args[1:], " "))
		return err
	case "kill":
		if len(args) < 2 {
			return fmt.Errorf("kill: missing pid")
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("kill: invalid pid: %v", err)
		}
		return syscall.Kill(pid, syscall.SIGTERM)
	case "ps":
		entries, err := os.ReadDir("/proc")
		if err != nil {
			return fmt.Errorf("ps: %v", err)
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if _, err := strconv.Atoi(e.Name()); err == nil {
				cmdline, _ := os.ReadFile(filepath.Join("/proc", e.Name(), "cmdline"))
				line := strings.ReplaceAll(string(cmdline), "\x00", " ")
				fmt.Fprintf(out, "%s %s\n", e.Name(), strings.TrimSpace(line))
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown builtin: %s", args[0])
	}
}
