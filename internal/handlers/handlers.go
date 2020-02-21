package handlers

import (
	"log"
	"ncaabballstats/internal/errors"
	"ncaabballstats/internal/helpers"
	"net/http"
	"os"
	"strconv"
)

type App struct {
	TeamHandler *TeamHandler
}

type TeamHandler struct {
	YearHandler *YearHandler
}

type YearHandler struct {
	PerGameHandler *PerGameHandler
}

type PerGameHandler struct {
}

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		h.TeamHandler.Handler(head).ServeHTTP(w, r)
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
			if _, err := strconv.Atoi(head); err != nil || len(head) > 4 {
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
		PerGameHandler: new(PerGameHandler),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
		// no parameters should be after this
		if r.URL.Path == "/" {
			switch head {
			case "pergame":
				h.PerGameHandler.Handler(school, year, head).ServeHTTP(w, r)
				return
			default:
				w.WriteHeader(http.StatusNotFound)
				w.Write(errors.NilEndpoint())
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(errors.NilEndpoint())
	})
}

func (h *PerGameHandler) Handler(school string, year string, stat string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failed := helpers.Scrape(school, year, stat)
		if failed != false {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errors.InternalErr())
		} else {
			dir, err := os.Getwd()
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(errors.InternalErr())
			} else {
				http.ServeFile(w, r, dir+string(os.PathSeparator)+school+year+".json")
			}
		}

	})
}
