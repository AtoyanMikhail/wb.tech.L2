package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type config struct {
	afterContext  int      
	beforeContext int      
	context       int      
	countOnly     bool     
	ignoreCase    bool     
	invertMatch   bool     
	fixedString   bool     
	lineNumber    bool     
	pattern       string   
	files         []string 
}

type matchResult struct {
	line     string 
	lineNum  int    
	filename string 
	matched  bool   
}


func compilePattern(cfg *config) (*regexp.Regexp, error) {
	pattern := cfg.pattern

	if cfg.fixedString {
		pattern = regexp.QuoteMeta(pattern)
	}

	if cfg.ignoreCase {
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}


func processFile(cfg *config, filename string, output io.Writer) (int, error) {
	var file io.ReadCloser
	var err error

	if filename == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return 0, err
		}
		defer file.Close()
	}

	re, err := compilePattern(cfg)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	results := []matchResult{}
	lineNum := 0
	matchCount := 0

	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := re.MatchString(line)

		if cfg.invertMatch {
			matched = !matched
		}

		if matched {
			matchCount++
		}

		results = append(results, matchResult{
			line:     line,
			lineNum:  lineNum,
			filename: filename,
			matched:  matched,
		})
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	
	if cfg.countOnly {
		fmt.Fprintf(output, "%d\n", matchCount)
		return matchCount, nil
	}

	
	printed := make(map[int]bool)
	totalLines := len(results)

	for i, res := range results {
		if res.matched {
			
			start := max(0, i-cfg.beforeContext)
			end := min(totalLines-1, i+cfg.afterContext)

			
			if cfg.context > 0 {
				start = max(0, i-cfg.context)
				end = min(totalLines-1, i+cfg.context)
			}

			
			for j := start; j <= end; j++ {
				if printed[j] {
					continue
				}
				printed[j] = true

				
				if cfg.lineNumber {
					fmt.Fprintf(output, "%d:", results[j].lineNum)
				}
				fmt.Fprintln(output, results[j].line)
			}
		}
	}

	return matchCount, nil
}

func main() {
	cfg := &config{}

	
	flag.IntVar(&cfg.afterContext, "A", 0, "Print N lines after match")
	flag.IntVar(&cfg.beforeContext, "B", 0, "Print N lines before match")
	flag.IntVar(&cfg.context, "C", 0, "Print N lines around match")
	flag.BoolVar(&cfg.countOnly, "c", false, "Print only count of matching lines")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&cfg.invertMatch, "v", false, "Invert match")
	flag.BoolVar(&cfg.fixedString, "F", false, "Interpret pattern as fixed string")
	flag.BoolVar(&cfg.lineNumber, "n", false, "Print line numbers")

	flag.Parse()

	
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: pattern required")
		os.Exit(1)
	}

	cfg.pattern = args[0]
	cfg.files = args[1:]

	
	if len(cfg.files) == 0 {
		_, err := processFile(cfg, "", os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			os.Exit(1)
		}
	} else {
		totalCount := 0
		for _, file := range cfg.files {
			count, err := processFile(cfg, file, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", file, err)
				continue
			}
			totalCount += count
		}

		
		if cfg.countOnly && len(cfg.files) > 1 {
			fmt.Printf("%d\n", totalCount)
		}
	}
}
