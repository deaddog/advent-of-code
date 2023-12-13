package main

import (
	"adventofcode/day1"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

var dayPathRegex = regexp.MustCompile("^/day([0-9]+)")

func main() {

	days := map[int]func() templ.Component{}
	days[1] = day1.DayTemplate

	slice := []int{}
	for d := range days {
		slice = append(slice, d)
	}

	http.Handle("/static/", NoCache(http.FileServer(http.Dir(""))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parsed := dayPathRegex.FindStringSubmatch(r.URL.Path)
		if len(parsed) == 2 {
			day, err := strconv.Atoi(parsed[1])
			if err != nil {
				r.Response.StatusCode = http.StatusInternalServerError
				return
			}

			if comp, ok := days[day]; ok {
				if r.Header.Get("Hx-Request") != "true" {
					templ.Handler(layout(slice, comp())).ServeHTTP(w, r)
				} else {
					templ.Handler(comp()).ServeHTTP(w, r)
				}
				return
			}
		}

		templ.Handler(layout(slice, root())).ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
