package handlers

import (
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"log"
	"ncaabballstats/internal/errors"
	"ncaabballstats/internal/helpers"
	"net/http"
	"os"
	"strconv"
)

type App struct {
	TeamHandler *TeamHandler
	ApiHandler  *ApiHandler
}

type ApiHandler struct {
}

type TeamHandler struct {
	YearHandler *YearHandler
}

type YearHandler struct {
	StatHandler *StatHandler
}

type StatHandler struct {
}

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		h.TeamHandler.Handler(head).ServeHTTP(w, r)
		return
	}
	if head == "api" {
		h.ApiHandler.ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(errors.NilEndpoint())
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, r.URL.Path = helpers.ShiftPath(r.URL.Path)
	// no parameters should be after this
	if r.URL.Path == "/" {
		file, err := ioutil.ReadFile("api/PerGame.md")
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errors.InternalErr())
		}
		api := github_flavored_markdown.Markdown(file)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(api)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(errors.NilEndpoint())
}

func (h *TeamHandler) Handler(school string) http.Handler {
	h = &TeamHandler{
		YearHandler: new(YearHandler),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
		if r.URL.Path != "/" {
			if _, err := strconv.Atoi(head); err != nil || len(head) != 4 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errors.InvalidYearErr())
				return
			} else {
				h.YearHandler.Handler(school, head).ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(errors.NilEndpoint())
	})
}

func (h *YearHandler) Handler(school string, year string) http.Handler {
	h = &YearHandler{
		StatHandler: new(StatHandler),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
		// no parameters should be after this
		if r.URL.Path == "/" {
			switch head {
			case "pergame":
				h.StatHandler.Handler(school, year, "per_game").ServeHTTP(w, r)
				return
			default:
				w.WriteHeader(http.StatusNotFound)
				w.Write(errors.NilEndpoint())
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(errors.NilEndpoint())
	})
}

func (h *StatHandler) Handler(school string, year string, stat string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failed := helpers.Scrape(school, year, stat)
		if failed == "internal" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errors.InternalErr())
			return
		}
		if failed == "noResponse" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(errors.ResponseErr())
			return
		}
		dir, err := os.Getwd()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errors.InternalErr())
			return
		}
		http.ServeFile(w, r, dir+string(os.PathSeparator)+school+year+stat+".json")
	})
}
