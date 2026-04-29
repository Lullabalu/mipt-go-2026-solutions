//go:build !solution

package httpgauge

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"sync"
)

type Gauge struct {
	mu sync.Mutex

	mp map[string]int
}

func New() *Gauge {
	return &Gauge{mp: make(map[string]int)}
}

func (g *Gauge) Snapshot() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()

	cp := make(map[string]int, len(g.mp))
	for k, v := range g.mp {
		cp[k] = v
	}
	return cp
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	snap := g.Snapshot()

	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s %d\n", k, snap[k])
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rc := chi.RouteContext(r.Context())
			pattern := rc.RoutePattern()
			if pattern == "" {
				return
			}
			g.mu.Lock()
			g.mp[pattern]++
			g.mu.Unlock()
		}()

		next.ServeHTTP(w, r)
	})
}
