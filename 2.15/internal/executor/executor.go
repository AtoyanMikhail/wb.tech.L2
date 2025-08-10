package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"wb-l2/internal/builtins"
	"wb-l2/internal/model"
)

// RunPipeline запускает пайплайн команд. Ctrl+C должен отменять контекст upstream.
func RunPipeline(ctx context.Context, cmds []model.Command) error {
	if len(cmds) == 1 && builtins.IsBuiltin(cmds[0].Args[0]) && cmds[0].InFile == "" && cmds[0].OutFile == "" {
		return builtins.Run(cmds[0].Args, os.Stdin, os.Stdout)
	}

	var procs []*exec.Cmd
	var readers []io.ReadCloser
	var writers []io.WriteCloser

	var in io.ReadCloser = os.Stdin
	var err error
	if cmds[0].InFile != "" {
		in, err = os.Open(cmds[0].InFile)
		if err != nil {
			return err
		}
		readers = append(readers, in)
	}

	for i, c := range cmds {
		var cmd *exec.Cmd
		if builtins.IsBuiltin(c.Args[0]) {
			cmd = builtinAsCmd(c.Args)
		} else {
			cmd = exec.CommandContext(ctx, c.Args[0], c.Args[1:]...)
		}

		if i == 0 {
			cmd.Stdin = in
		}
		if i < len(cmds)-1 {
			pr, pw := io.Pipe()
			cmd.Stdout = pw
			readers = append(readers, pr)
			writers = append(writers, pw)
		} else {
			if c.OutFile != "" {
				var f *os.File
				if c.AppendOut {
					f, err = os.OpenFile(c.OutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
				} else {
					f, err = os.OpenFile(c.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
				}
				if err != nil {
					return err
				}
				writers = append(writers, f)
				cmd.Stdout = f
			} else {
				cmd.Stdout = os.Stdout
			}
		}
		cmd.Stderr = os.Stderr
		procs = append(procs, cmd)
	}

	var prevReader io.ReadCloser
	if cmds[0].InFile != "" {
		prevReader = readers[0]
	}
	rIdx := 0
	for i := range procs {
		if i > 0 {
			if prevReader == nil {
				prevReader = readers[rIdx]
			}
			procs[i].Stdin = prevReader
		}
		if i < len(procs)-1 {
			if cmds[0].InFile != "" {
				rIdx = i
				prevReader = readers[rIdx+1]
			} else {
				rIdx = i
				prevReader = readers[rIdx]
			}
		}
	}

	for i, p := range procs {
		if err := p.Start(); err != nil {
			for _, r := range readers {
				_ = r.Close()
			}
			for _, w := range writers {
				_ = w.Close()
			}
			return err
		}
		if i > 0 && i-1 < len(writers) {
			_ = writers[i-1].Close()
		}
	}

	done := make(chan error, len(procs))
	for _, p := range procs {
		go func(cmd *exec.Cmd) { done <- cmd.Wait() }(p)
	}

	timeout := time.NewTimer(0)
	if !timeout.Stop() {
		<-timeout.C
	}

	var firstErr error
	for i := 0; i < len(procs); i++ {
		select {
		case err := <-done:
			if err != nil && firstErr == nil {
				firstErr = err
			}
		case <-ctx.Done():
			for _, p := range procs {
				if p.Process != nil {
					_ = p.Process.Kill()
				}
			}
			firstErr = ctx.Err()
		}
	}

	for _, r := range readers {
		_ = r.Close()
	}
	for _, w := range writers {
		_ = w.Close()
	}
	return firstErr
}

func builtinAsCmd(args []string) *exec.Cmd {
	switch args[0] {
	case "echo":
		return exec.Command("/bin/sh", "-c", "echo "+shellQuote(strings.Join(args[1:], " ")))
	case "pwd":
		return exec.Command("/bin/sh", "-c", "pwd")
	case "ps":
		return exec.Command(os.Args[0], "__builtin_ps")
	case "kill":
		if len(args) < 2 {
			return exec.Command("/bin/sh", "-c", "false")
		}
		return exec.Command(os.Args[0], "__builtin_kill", args[1])
	default:
		return exec.Command("/bin/sh", "-c", "false")
	}
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if strings.IndexFunc(s, func(r rune) bool { return r == ' ' || r == '\'' || r == '"' }) >= 0 {
		return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
	}
	return s
}
