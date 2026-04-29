//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func MakeHandler(service any) http.Handler {
	sv := reflect.ValueOf(service)

	ctxT := reflect.TypeOf((*context.Context)(nil)).Elem()
	errT := reflect.TypeOf((*error)(nil)).Elem()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.Trim(r.URL.Path, "/")
		if name == "" {
			http.NotFound(w, r)
			return
		}

		m := sv.MethodByName(name)
		if !m.IsValid() {
			http.NotFound(w, r)
			return
		}

		mt := m.Type()
		if mt.NumIn() != 2 || !mt.In(0).Implements(ctxT) || mt.In(1).Kind() != reflect.Ptr ||
			mt.NumOut() != 2 || mt.Out(0).Kind() != reflect.Ptr || !mt.Out(1).Implements(errT) {
			http.NotFound(w, r)
			return
		}

		reqPtr := reflect.New(mt.In(1).Elem())
		if err := json.NewDecoder(r.Body).Decode(reqPtr.Interface()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out := m.Call([]reflect.Value{reflect.ValueOf(r.Context()), reqPtr})
		if !out[1].IsNil() {
			http.Error(w, out[1].Interface().(error).Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out[0].Interface())
	})
}
func Call(ctx context.Context, endpoint string, method string, req, rsp any) error {
	url := endpoint + "/" + method
	bts, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	json.NewDecoder(resp.Body).Decode(rsp)
	return nil
}
