package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"wb-l2/2.16/internal/htmlparse"
	"wb-l2/2.16/internal/store"
	"wb-l2/2.16/internal/urlutil"
)

type Config struct {
	StartURL    *url.URL
	MaxDepth    int
	OutputDir   string
	Concurrency int
	Timeout     time.Duration
}

type task struct {
	url   *url.URL
	depth int
}

func Run(ctx context.Context, cfg Config) error {
	client := &http.Client{Timeout: cfg.Timeout}
	visited := sync.Map{}
	sem := make(chan struct{}, cfg.Concurrency)
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)

	var submit func(u *url.URL, depth int)
	submit = func(u *url.URL, depth int) {
		if depth > cfg.MaxDepth {
			return
		}
		key := u.String()
		if _, loaded := visited.LoadOrStore(key, struct{}{}); loaded {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()
			if err := fetchOne(ctx, client, cfg.OutputDir, cfg.StartURL, u, depth, submit); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	submit(cfg.StartURL, 0)
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func fetchOne(ctx context.Context, client *http.Client, outDir string, root *url.URL, u *url.URL, depth int, submit func(*url.URL, int)) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", u, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ct := resp.Header.Get("Content-Type")
	page, _ := htmlparse.Extract(u, ct, body)

	// compute local path
	local := urlutil.LocalPath(outDir, u)

	// For HTML: переписываем ссылки на локальные пути
	if page.IsHTML && len(page.Links) > 0 {
		repl := make(map[string]string)
		for _, ref := range page.Links {
			if nu, ok := urlutil.Normalize(u, ref); ok {
				repl[ref] = urlutil.LocalPath(outDir, nu)
			}
		}
		if len(repl) > 0 && page.Rewrite != nil {
			body = page.Rewrite(repl)
		}
	}
	if err := store.Save(local, body); err != nil {
		return err
	}

	// Submit links (only same-host)
	for _, ref := range page.Links {
		if nu, ok := urlutil.Normalize(u, ref); ok {
			submit(nu, depth+1)
		}
	}
	_ = root
	return nil
}

func pageIsHTML(ct string) bool { return strings.Contains(strings.ToLower(ct), "text/html") }
