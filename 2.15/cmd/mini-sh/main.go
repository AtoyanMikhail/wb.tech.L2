package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"wb-l2/internal/builtins"
	"wb-l2/internal/executor"
	"wb-l2/internal/parser"
)

func isTTY(fd uintptr) bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	reader := bufio.NewReader(os.Stdin)
	for {
		if isTTY(os.Stdin.Fd()) {
			fmt.Print(parser.Prompt())
		}

		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = os.ExpandEnv(line)

		orSegments := parser.SplitPreserve(line, "||")
		var lastOK bool
		for idx, seg := range orSegments {
			andSegments := parser.SplitPreserve(seg, "&&")
			lastOK = true
			for _, a := range andSegments {
				a = strings.TrimSpace(a)
				if a == "" {
					continue
				}

				cmds, perr := parser.ParsePipeline(a)
				if perr != nil {
					fmt.Fprintln(os.Stderr, perr)
					lastOK = false
					break
				}

				ctx, cancel := context.WithCancel(context.Background())
				done := make(chan struct{})
				go func() {
					select {
					case <-sigCh:
						cancel()
					case <-done:
					}
				}()

				err = executor.RunPipeline(ctx, cmds)
				close(done)
				cancel()
				if err != nil {
					lastOK = false
					break
				}
			}
			if lastOK {
				break
			}
			if idx == len(orSegments)-1 {
				// last OR segment failed; loop next prompt
			}
		}
		_ = lastOK
	}
}

// builtin helpers
func init() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "__builtin_ps":
			_ = builtins.Run([]string{"ps"}, os.Stdin, os.Stdout)
			os.Exit(0)
		case "__builtin_kill":
			if len(os.Args) >= 3 {
				_ = builtins.Run([]string{"kill", os.Args[2]}, os.Stdin, os.Stdout)
			}
			os.Exit(0)
		}
	}
}
