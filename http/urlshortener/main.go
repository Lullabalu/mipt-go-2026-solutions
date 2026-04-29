//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func main() {
	mp := make(map[string]string)
	mpu := make(map[string]string)
	// var mu sync.Mutex
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	handlerPost := func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key := ""
		if _, ok := mpu[req.URL]; ok {
			key = mpu[req.URL]
		} else {
			key = strconv.Itoa(len(mp))
		}

		mp[key] = req.URL
		mpu[req.URL] = key
		resp := Response{URL: req.URL, Key: key}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}

	handlerGet := func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/go/")
		url, ok := mp[key]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	}

	http.HandleFunc("/shorten", handlerPost)
	http.HandleFunc("/go/", handlerGet)
	http.ListenAndServe(":"+*port, nil)

}
