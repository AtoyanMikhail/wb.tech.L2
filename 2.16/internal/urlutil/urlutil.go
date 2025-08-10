package urlutil

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// Normalize resolves ref relative to base and strips fragments
func Normalize(base *url.URL, ref string) (*url.URL, bool) {
	r, err := url.Parse(ref)
	if err != nil {
		return nil, false
	}
	u := base.ResolveReference(r)
	u.Fragment = ""
	if u.Scheme != base.Scheme || u.Host != base.Host {
		return nil, false
	}
	return u, true
}

// LocalPath converts URL to a filesystem path within outDir
func LocalPath(outDir string, u *url.URL) string {
	p := u.Path
	if p == "" || strings.HasSuffix(p, "/") {
		p = path.Join(p, "index.html")
	}
	if !strings.Contains(path.Base(p), ".") {
		// no extension: treat as html page
		p = path.Join(p, "index.html")
	}
	return filepath.Join(outDir, u.Host, filepath.FromSlash(p))
}
