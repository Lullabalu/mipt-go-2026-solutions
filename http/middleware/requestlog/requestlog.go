//go:build !solution

package requestlog

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Id struct {
	id int

	mu sync.Mutex
}

type MyWriter struct {
	http.ResponseWriter
	statusCode int
}

func (wr *MyWriter) WriteHeader(code int) {
	wr.statusCode = code
	wr.ResponseWriter.WriteHeader(code)
}

func (i *Id) GetId() (id int) {
	i.mu.Lock()
	id = i.id
	i.id += 1
	i.mu.Unlock()
	return
}

var id Id

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			curId := id.GetId()
			m := r.Method
			p := r.URL.Path
			curTime := time.Now()

			mw := &MyWriter{ResponseWriter: w, statusCode: 200}

			l.Info("request started",
				zap.String("path", p),
				zap.String("method", m),
				zap.Int("request_id", curId),
			)
			defer func() {
				duration := time.Since(curTime)
				msg := "request finished"
				panic_msg := ""
				if v := recover(); v != nil {
					msg = "request panicked"
					panic_msg = fmt.Sprintf("%v", v)
				}
				l.Info(msg,
					zap.String("path", p),
					zap.String("method", m),
					zap.Int("request_id", curId),
					zap.Int("status_code", mw.statusCode),
					zap.Duration("duration", duration),
				)
				if msg == "request panicked" {
					panic(panic_msg)
				}
			}()
			next.ServeHTTP(mw, r)

		})
	}
}
