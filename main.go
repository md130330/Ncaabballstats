package main

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"web_scraper/excel"
)

type ResponseError struct {
	StatusCode int
	RequestUrl *url.URL
}

func perGameStats(w http.ResponseWriter, r *http.Request) {
	school := mux.Vars(r)["school"]
	year := mux.Vars(r)["year"]

	c := colly.NewCollector()

	c.OnHTML("div[id=div_per_game]", func(e *colly.HTMLElement) {

		dir, err := os.Getwd()
		if err != nil {
			w.Write([]byte(`{"message": "Internal Error. Please contact administrator for more details"}`))
			log.Fatal(err)
		}

		e.DOM.Find("col").Remove()
		e.DOM.Find("caption").Remove()
		table, err := e.DOM.Html()
		if err != nil {
			w.Write([]byte(`{"message": "Internal Error. Please contact administrator for more details"}`))
			log.Fatal(err)
		}
		byteTable := []byte(table)
		writeErr := ioutil.WriteFile(dir+"\\"+school+year+".xls", byteTable, 0644)
		if writeErr != nil {
			w.Write([]byte(`{"message": "Internal Error. Please contact administrator for more details"}`))
			log.Fatal(err)
		}

		xlsFile := &excel.Excel{Visible: false, Readonly: false, Saved: true}
		filePath := dir + "\\" + school + year + ".xls"
		xlsFile.Open(filePath)
		xlsFile.SaveAs(dir+"\\"+school+year, "csv")
		xlsFile.Close()

		cmd := exec.Command("./convert_to_json", "-path="+dir+"\\"+school+year+".csv")
		cmd.Run()
		http.ServeFile(w, r, dir+"\\"+school+year+".json")
	})

	c.OnError(func(r *colly.Response, rErr error) {
		responseErr := ResponseError{r.StatusCode, r.Request.URL}
		errJson, err := json.Marshal(responseErr)
		if err != nil {

			log.Fatal(err)
		}
		w.Write(errJson)
	})

	c.Visit("https://www.sports-reference.com/cbb/schools/" + school + "/" + year + ".html")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods(http.MethodGet)
	r.HandleFunc("/{school}/{year}", perGameStats).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
