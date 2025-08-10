package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wb-l2/internal/model"
)

// Prompt возвращает строку приглашения
func Prompt() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("mini-sh:%s$ ", filepath.Base(cwd))
}

// SplitPreserve делит строку по сепаратору без регэкспов
func SplitPreserve(s, sep string) []string {
	if !strings.Contains(s, sep) {
		return []string{s}
	}
	var parts []string
	cur := s
	for {
		idx := strings.Index(cur, sep)
		if idx < 0 {
			parts = append(parts, cur)
			break
		}
		parts = append(parts, cur[:idx])
		cur = cur[idx+len(sep):]
	}
	return parts
}

// ParsePipeline парсит пайплайн команды в список стадий
func ParsePipeline(line string) ([]model.Command, error) {
	stages := SplitPreserve(line, "|")
	cmds := make([]model.Command, 0, len(stages))
	for _, st := range stages {
		st = strings.TrimSpace(st)
		if st == "" {
			return nil, fmt.Errorf("empty stage in pipeline")
		}
		cmd, err := parseRedirections(st)
		if err != nil {
			return nil, err
		}
		if len(cmd.Args) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

// parseRedirections обрабатывает редиректы <, >, >>
func parseRedirections(s string) (model.Command, error) {
	var cmd model.Command
	tokens := strings.Fields(s)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t {
		case ">":
			if i+1 >= len(tokens) {
				return cmd, fmt.Errorf("redirect > requires filename")
			}
			cmd.OutFile = tokens[i+1]
			cmd.AppendOut = false
			i++
		case ">>":
			if i+1 >= len(tokens) {
				return cmd, fmt.Errorf("redirect >> requires filename")
			}
			cmd.OutFile = tokens[i+1]
			cmd.AppendOut = true
			i++
		case "<":
			if i+1 >= len(tokens) {
				return cmd, fmt.Errorf("redirect < requires filename")
			}
			cmd.InFile = tokens[i+1]
			i++
		default:
			cmd.Args = append(cmd.Args, t)
		}
	}
	return cmd, nil
}
