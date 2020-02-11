package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gocolly/colly/v2"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"net/url"
	"os"
	"strings"
	"web_scraper/csv-to-json"
)

type ResponseError struct {
	StatusCode int
	Response   *url.URL
}

type InternalError struct {
	StatusCode int
	Response   string
}

const ErrorMsg = "Internal Error. Please contact administrator for more details"

func perGameStats(w http.ResponseWriter, r *http.Request) {
	school := mux.Vars(r)["school"]
	year := mux.Vars(r)["year"]

	c := colly.NewCollector()

	c.OnHTML("div[id=div_per_game]", func(e *colly.HTMLElement) {

		p := bluemonday.StrictPolicy()

		var table [][]string
		var header []string

		dir, err := os.Getwd()
		if err != nil {
			log.Info(err)
			internalErr := InternalError{500, ErrorMsg}
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
				internalErr := InternalError{500, ErrorMsg}
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
			internalErr := InternalError{500, ErrorMsg}
			errJson, _ := json.Marshal(internalErr)
			w.Write(errJson)
		}

		writer := csv.NewWriter(file)

		for _, value := range table {
			err := writer.Write(value)
			if err != nil {
				log.Info(err)
				internalErr := InternalError{500, ErrorMsg}
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
			internalErr := InternalError{500, ErrorMsg}
			errJson, _ := json.Marshal(internalErr)
			w.Write(errJson)
		} else {
			http.ServeFile(w, r, dir+string(os.PathSeparator)+school+year+".json")
		}
	})

	c.OnError(func(r *colly.Response, rErr error) {
		responseErr := ResponseError{r.StatusCode, r.Request.URL}
		errJson, _ := json.Marshal(responseErr)
		w.Write(errJson)
	})

	c.Visit("https://www.sports-reference.com/cbb/schools/" + school + "/" + year + ".html")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", home).Methods(http.MethodGet)

	r.HandleFunc("/{school}/{year}/pergame", perGameStats).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
