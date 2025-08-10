package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type key struct {
	strVal   string
	numVal   float64
	isNum    bool
	original string
}

type application struct {
	column               int
	numeric              bool
	reverse              bool
	unique               bool
	sortByMonth          bool
	ignoreTrailingBlanks bool
	checkSort            bool
	humanNumeric         bool
}

func (app *application) key(s string) *key {
	field := s

	if app.column > 0 {
		fields := strings.Split(s, "\t")
		idx := app.column - 1
		if idx < len(fields) {
			field = fields[idx]
		} else {
			field = ""
		}
	}

	if app.ignoreTrailingBlanks {
		field = strings.TrimRight(field, " \t")
	}

	key := &key{original: s, strVal: field}

	switch {
	case app.sortByMonth:
		if month, ok := parseMonth(field); ok {
			key.isNum = true
			key.numVal = float64(month)
		}
	case app.humanNumeric:
		if num, ok := parseHumanNumber(field); ok {
			key.isNum = true
			key.numVal = num
		}
	case app.numeric:
		if num, ok := parseNumber(field); ok {
			key.isNum = true
			key.numVal = num
		}
	}

	return key
}

func keyLess(a, b *key) bool {
	if a.isNum {
		if b.isNum {
			return a.numVal < b.numVal
		}
		return false
	}
	if b.isNum {
		return true
	}
	return a.strVal < b.strVal
}

func (app *application) fullLess(a, b *key) bool {

	if keyLess(a, b) {
		return true
	}
	if keyLess(b, a) {
		return false
	}

	return a.original < b.original
}

func parseMonth(s string) (int, bool) {
	months := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	s = strings.TrimSpace(s)
	for i, month := range months {
		if strings.EqualFold(s, month) {
			return i + 1, true
		}
	}
	return 0, false
}

func parseHumanNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	var numStr, suffix string
	for i, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			numStr += string(r)
		} else {
			suffix = strings.ToUpper(strings.TrimSpace(s[i:]))
			break
		}
	}

	multipliers := map[string]float64{
		"K": 1e3, "M": 1e6, "G": 1e9, "T": 1e12, "P": 1e15, "E": 1e18,
		"К": 1e3, "М": 1e6, "Г": 1e9, "Т": 1e12,
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, false
	}

	if mult, ok := multipliers[suffix]; ok {
		return num * mult, true
	}
	return num, true
}

func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return num, true
}

func readLines(filenames []string) ([]string, error) {
	if len(filenames) == 0 {
		return readFromStdin()
	}

	var lines []string
	for _, filename := range filenames {
		fileLines, err := readFromFile(filename)
		if err != nil {
			return nil, err
		}
		lines = append(lines, fileLines...)
	}
	return lines, nil
}

func readFromStdin() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func readFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (app *application) checkSorted(lines []string) bool {
	var prev *key
	for i, line := range lines {
		key := app.key(line)
		if i > 0 {
			if app.reverse {
				if keyLess(prev, key) {
					fmt.Fprintf(os.Stderr, "disorder: %s\n", line)
					return false
				}
			} else {
				if keyLess(key, prev) {
					fmt.Fprintf(os.Stderr, "disorder: %s\n", line)
					return false
				}
			}
		}
		prev = key
	}
	return true
}

func removeDuplicates(keys []*key, lines []string) ([]*key, []string) {
	if len(keys) == 0 {
		return keys, lines
	}

	i := 0
	for j := 1; j < len(keys); j++ {
		if keys[i].isNum != keys[j].isNum ||
			(keys[i].isNum && keys[i].numVal != keys[j].numVal) ||
			(!keys[i].isNum && keys[i].strVal != keys[j].strVal) {
			i++
			keys[i] = keys[j]
			lines[i] = lines[j]
		}
	}
	return keys[:i+1], lines[:i+1]
}
