package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"wb-l2/2.16/internal/crawler"
)

func main() {
	var depth int
	var outDir string
	var concurrency int
	flag.IntVar(&depth, "depth", 2, "recursion depth")
	flag.StringVar(&outDir, "o", "site", "output directory")
	flag.IntVar(&concurrency, "j", 6, "max concurrent downloads")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: wget [-depth N] [-o DIR] [-j N] <url>\n")
		os.Exit(2)
	}
	raw := flag.Arg(0)
	u, err := url.Parse(raw)
	if err != nil || !u.IsAbs() || u.Host == "" {
		log.Fatalf("invalid url: %q", raw)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg := crawler.Config{
		StartURL:    u,
		MaxDepth:    depth,
		OutputDir:   outDir,
		Concurrency: concurrency,
		Timeout:     15 * time.Second,
	}
	if err := crawler.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
