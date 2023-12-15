package main

import (
	"adventofcode/day1"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

var dayPathRegex = regexp.MustCompile("^/day([0-9]+)")

type DayConfig struct {
	Filename  string
	Component func(string) (templ.Component, error)
}

func main() {
	days := map[int]DayConfig{
		1: {
			Filename:  "day1",
			Component: day1.Solution,
		},
	}

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

			if config, ok := days[day]; ok {
				component := config.getComponent()

				if r.Header.Get("Hx-Request") != "true" {
					component = layout(slice, component)
				}

				templ.Handler(component).ServeHTTP(w, r)
				return
			}
		}

		templ.Handler(layout(slice, root())).ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (cfg *DayConfig) getData() (string, error) {
	p := os.Getenv("ADVENTOFCODE_DATA")
	if p == "" {
		return "", fmt.Errorf("no data path set")
	}

	data, err := os.ReadFile(path.Join(p, cfg.Filename))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (cfg *DayConfig) getComponent() templ.Component {
	content, err := cfg.getData()
	if err != nil {
		return ErrTemplate(err)
	}

	component, err := cfg.Component(content)
	if err != nil {
		return ErrTemplate(err)
	}

	return component
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
