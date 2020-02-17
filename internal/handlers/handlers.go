package handlers

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	log "github.com/golang/glog"
	"github.com/microcosm-cc/bluemonday"
	"ncaabballstats/internal/errors"
	"ncaabballstats/internal/helpers"
	"ncaabballstats/pkg/csv-to-json"
	"net/http"
	"os"
	"strings"
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
	http.Error(w, "Not Found app", http.StatusNotFound)
}

func (h *TeamHandler) Handler(school string) http.Handler {
	h = &TeamHandler{
		YearHandler: new(YearHandler),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
		if r.URL.Path != "/" {
			h.YearHandler.Handler(school, head).ServeHTTP(w, r)
			return
		}
		http.Error(w, "Not Found team", http.StatusNotFound)
	})
}

func (h *YearHandler) Handler(school string, year string) http.Handler {
	h = &YearHandler{
		PerGameHandler: new(PerGameHandler),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = helpers.ShiftPath(r.URL.Path)
		switch head {
		case "pergame":
			h.PerGameHandler.Handler(school, year, head).ServeHTTP(w, r)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})
}

func (h *PerGameHandler) Handler(school string, year string, stat string) http.Handler {
	h = &PerGameHandler{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := colly.NewCollector()

		internalErr := errors.InternalError{
			StatusCode: 500,
			Response:   errors.ErrorMsg,
		}

		c.OnHTML("div[id=div_per_game]", func(e *colly.HTMLElement) {

			p := bluemonday.StrictPolicy()

			var table [][]string
			var header []string

			dir, err := os.Getwd()
			if err != nil {
				log.Info(err)
				errJson, _ := json.Marshal(internalErr)
				w.Write(errJson)
			}

			e.DOM.Find("colgroup").Remove()
			e.DOM.Find("caption").Remove()
			e.ForEach("thead th", func(_ int, el *colly.HTMLElement) {
				header = append(header, el.Attr("aria-label"))
			})
			table = append(table, header)

			e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
				f, err := el.DOM.Html()
				if err != nil {
					log.Info(err)
					errJson, _ := json.Marshal(internalErr)
					w.Write(errJson)
				}
				f = strings.Replace(f, "</t", ",</t", -1)
				f = p.Sanitize(f)
				row := strings.Split(f, ",")
				row = row[:len(row)-1]
				for i := range row {
					if strings.HasPrefix(row[i], ".") {
						row[i] = strings.Replace(row[i], ".", "0.", -1)
					}
				}
				table = append(table, row)
			})
			file, err := os.Create(school + year + ".csv")
			if err != nil {
				log.Info(err)
				errJson, _ := json.Marshal(internalErr)
				w.Write(errJson)
			}

			writer := csv.NewWriter(file)

			for _, value := range table {
				err := writer.Write(value)
				if err != nil {
					log.Info(err)
					errJson, _ := json.Marshal(internalErr)
					w.Write(errJson)
				}
				writer.Flush()
			}

			file.Close()

			s := dir + string(os.PathSeparator) + school + year + ".csv"
			sPointer := &s
			_, convertErr := converter.Convert(sPointer)
			if convertErr != nil {
				log.Info(convertErr)
				errJson, _ := json.Marshal(internalErr)
				w.Write(errJson)
			} else {
				http.ServeFile(w, r, dir+string(os.PathSeparator)+school+year+".json")
			}
		})

		c.OnError(func(r *colly.Response, rErr error) {
			responseErr := errors.ResponseError{r.StatusCode, r.Request.URL}
			errJson, _ := json.Marshal(responseErr)
			w.Write(errJson)
		})

		c.Visit("https://www.sports-reference.com/cbb/schools/" + school + "/" + year + ".html")
	})
}
