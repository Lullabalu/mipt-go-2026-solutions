//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type AInfo struct {
	Name    string `json:"athlete"`
	sport   string
	Country string         `json:"country"`
	Medals  map[string]int `json:"medals"`

	MedalsByYear map[int]map[string]int `json:"medals_by_year"`
}

type CInfo struct {
	Name string `json:"country"`

	results map[int]map[string]int
}

type Data struct {
	Name    string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type Ans struct {
	Country string `json:"country"`

	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

func main() {
	port := flag.String("port", "8080", "server port")
	path := flag.String("data", "./testdata/olympicWinners.json", "path to json")
	flag.Parse()

	bytes, _ := os.ReadFile(*path)

	var rows []Data
	json.Unmarshal(bytes, &rows)

	athlets := make(map[string]AInfo)
	sportAth := make(map[string]map[string]*AInfo)
	countries := make(map[string]CInfo)

	for _, row := range rows {
		m, ok := sportAth[row.Sport]
		if !ok {
			m = make(map[string]*AInfo)
			sportAth[row.Sport] = m
		}

		ai, ok := m[row.Name]
		if !ok {

			ai = &AInfo{
				Name:         row.Name,
				Country:      row.Country,
				sport:        row.Sport,
				Medals:       make(map[string]int),
				MedalsByYear: make(map[int]map[string]int),
			}
			m[row.Name] = ai
		}

		ai.Medals["gold"] += row.Gold
		ai.Medals["silver"] += row.Silver
		ai.Medals["bronze"] += row.Bronze
		ai.Medals["total"] += row.Total

		ym, ok := ai.MedalsByYear[row.Year]
		if !ok {
			ym = make(map[string]int)
			ai.MedalsByYear[row.Year] = ym
		}
		ym["gold"] += row.Gold
		ym["silver"] += row.Silver
		ym["bronze"] += row.Bronze
		ym["total"] += row.Total

		cn, ok := countries[row.Country]
		if !ok {
			cn = CInfo{
				Name:    row.Country,
				results: make(map[int]map[string]int),
			}
		}
		athlets[row.Name] = *ai

		curResults, ok := cn.results[row.Year]
		if !ok {
			curResults = make(map[string]int)
		}
		curResults["gold"] += row.Gold
		curResults["silver"] += row.Silver
		curResults["bronze"] += row.Bronze
		curResults["total"] += row.Total

		cn.results[row.Year] = curResults

		countries[row.Country] = cn
	}
	http.Handle("/athlete-info", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := r.URL.Query().Get("name")

			if name == "" {
				http.Error(w, "missing query param: name", http.StatusBadRequest)
				return
			}

			ai, ok := athlets[name]
			if !ok {
				http.Error(w, "athlete not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(ai)
		}))

	http.Handle("/top-athletes-in-sport", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sport := r.URL.Query().Get("sport")

			if sport == "" {
				http.Error(w, "spirt", http.StatusBadRequest)
				return
			}

			limits := r.URL.Query().Get("limit")

			if limits == "" {
				limits = "3"
			}

			limit, err := strconv.Atoi(limits)
			if err != nil || limit <= 0 {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}

			m := sportAth[sport]
			if m == nil {
				http.Error(w, "sport not found", http.StatusNotFound)
				return
			}

			ais := make([]AInfo, 0, len(m))
			for _, p := range m {
				ais = append(ais, *p)
			}

			if len(ais) == 0 {
				http.Error(w, "sport not found", http.StatusNotFound)
				return
			}

			sort.Slice(ais, func(i, j int) bool {
				switch {
				case ais[i].Medals["gold"] != ais[j].Medals["gold"]:
					return ais[i].Medals["gold"] > ais[j].Medals["gold"]
				case ais[i].Medals["silver"] != ais[j].Medals["silver"]:
					return ais[i].Medals["silver"] > ais[j].Medals["silver"]
				case ais[i].Medals["bronze"] != ais[j].Medals["bronze"]:
					return ais[i].Medals["bronze"] > ais[j].Medals["bronze"]
				default:
					return ais[i].Name < ais[j].Name
				}
			})

			if limit > len(ais) {
				limit = len(ais)
			}
			ais = ais[:limit]

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(ais)

		}))

	http.Handle("/top-countries-in-year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			years := r.URL.Query().Get("year")
			if years == "" {
				http.Error(w, "spirt", http.StatusBadRequest)
				return
			}

			year, err := strconv.Atoi(years)
			if err != nil {
				http.Error(w, "invalid year", http.StatusBadRequest)
				return
			}
			limits := r.URL.Query().Get("limit")
			if limits == "" {
				limits = "3"
			}
			limit, err := strconv.Atoi(limits)
			if err != nil || limit <= 0 {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}

			cns := []Ans{}
			for _, cn := range countries {
				results, ok := cn.results[year]
				if ok {
					cns = append(cns, Ans{Country: cn.Name, Gold: results["gold"], Silver: results["silver"], Bronze: results["bronze"],
						Total: results["total"]})
				}
			}
			if len(cns) == 0 {
				http.Error(w, "sport not found", http.StatusNotFound)
				return
			}

			sort.Slice(cns, func(i, j int) bool {
				switch {
				case cns[i].Gold != cns[j].Gold:
					return cns[i].Gold > cns[j].Gold
				case cns[i].Silver != cns[j].Silver:
					return cns[i].Silver > cns[j].Silver
				case cns[i].Bronze != cns[j].Bronze:
					return cns[i].Bronze > cns[j].Bronze
				default:
					return cns[i].Country < cns[j].Country

				}
			})

			if limit > len(cns) {
				limit = len(cns)
			}
			cns = cns[:limit]

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(cns)

		}))
	http.ListenAndServe(":"+*port, nil)
}
