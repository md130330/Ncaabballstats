package helpers

import (
	"encoding/csv"
	"github.com/gocolly/colly/v2"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"ncaabballstats/pkg/csv-to-json"
	"os"
	"path"
	"strings"
)

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func Scrape(school string, year string, stat string) bool {
	var failed bool

	c := colly.NewCollector()

	log.Print("div[id=div_" + stat + "]")

	c.OnHTML("div[id=div_"+stat+"]", func(e *colly.HTMLElement) {

		p := bluemonday.StrictPolicy()

		var table [][]string
		var header []string

		dir, err := os.Getwd()
		if err != nil {
			log.Print(err)
			failed = true
		}

		e.DOM.Find("colgroup").Remove()
		e.DOM.Find("caption").Remove()
		e.ForEach("thead th", func(_ int, el *colly.HTMLElement) {
			header = append(header, el.Attr("aria-label"))
		})
		table = append(table, header)

		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
			f, err := el.DOM.Html()
			log.Print(f)
			if err != nil {
				log.Print(err)
				failed = true
			}
			f = strings.Replace(f, "</t", "~</t", -1)
			f = p.Sanitize(f)
			row := strings.Split(f, "~")
			row = row[:len(row)-1]
			for i := range row {
				if strings.HasPrefix(row[i], ".") {
					row[i] = strings.Replace(row[i], ".", "0.", -1)
				}
			}
			table = append(table, row)
		})
		file, err := os.Create(school + year + stat + ".csv")
		if err != nil {
			log.Print(err)
			failed = true
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		for _, value := range table {
			err := writer.Write(value)
			if err != nil {
				log.Print(err)
				failed = true
			}
			writer.Flush()
		}

		s := dir + string(os.PathSeparator) + school + year + stat + ".csv"
		sPointer := &s
		_, convertErr := converter.Convert(sPointer)
		if convertErr != nil {
			log.Print(convertErr)
			failed = true
		}
	})

	c.OnError(func(r *colly.Response, rErr error) {
		log.Print(rErr)
		failed = true
	})

	c.Visit("https://www.sports-reference.com/cbb/schools/" + school + "/" + year + ".html")

	return failed
}
