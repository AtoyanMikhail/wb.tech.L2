package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type config struct {
	fields    []int  
	delimiter string 
	separated bool   
}

func parseFields(fieldsStr string) ([]int, error) {
	var fields []int
	seen := make(map[int]bool) 

	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			
			rangeParts := strings.SplitN(part, "-", 2)
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil || start < 1 {
				return nil, fmt.Errorf("invalid start of range: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil || end < 1 {
				return nil, fmt.Errorf("invalid end of range: %s", rangeParts[1])
			}

			if start > end {
				start, end = end, start
			}

			for i := start; i <= end; i++ {
				idx := i - 1 
				if !seen[idx] {
					fields = append(fields, idx)
					seen[idx] = true
				}
			}
		} else {
			
			num, err := strconv.Atoi(part)
			if err != nil || num < 1 {
				return nil, fmt.Errorf("invalid field: %s", part)
			}

			idx := num - 1 
			if !seen[idx] {
				fields = append(fields, idx)
				seen[idx] = true
			}
		}
	}

	return fields, nil
}

func processLine(line string, cfg *config) string {
	if cfg.separated && !strings.Contains(line, cfg.delimiter) {
		return ""
	}

	
	parts := strings.Split(line, cfg.delimiter)

	
	if len(cfg.fields) == 0 {
		return ""
	}

	
	var outputParts []string
	for _, idx := range cfg.fields {
		if idx < len(parts) {
			outputParts = append(outputParts, parts[idx])
		}
	}

	
	return strings.Join(outputParts, cfg.delimiter)
}

func main() {
	fieldsFlag := flag.String("f", "", "Comma-separated list of fields to extract")
	delimiterFlag := flag.String("d", "\t", "Field delimiter character")
	separatedFlag := flag.Bool("s", false, "Suppress lines without delimiters")
	flag.Parse()
	
	if *fieldsFlag == "" {
		fmt.Fprintln(os.Stderr, "cut: you must specify a list of fields with -f")
		os.Exit(1)
	}
	
	fields, err := parseFields(*fieldsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cut: %v\n", err)
		os.Exit(1)
	}
	
	sort.Ints(fields)

	cfg := &config{
		fields:    fields,
		delimiter: *delimiterFlag,
		separated: *separatedFlag,
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		output := processLine(line, cfg)
		if output != "" {
			fmt.Fprintln(writer, output)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "cut: error reading input: %v\n", err)
		os.Exit(1)
	}
}