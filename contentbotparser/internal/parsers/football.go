package parsers

import (
	"github.com/gocolly/colly"
)

func ChampionatComParser() {
	url := "https://www.championat.com/"
	c := colly.NewCollector()

	c.OnHTML("div.news-items", func(e *colly.HTMLElement) {
		proccess := true
		e.ForEach("div", func(i int, element *colly.HTMLElement) {
			if proccess {
				if element.Attr("class") == "news-items__head" {
					if element.Text == "5 октября 2023" {
						proccess = false
					}
				}

			}
		})
	})

	c.Visit(url + "news/football/1.html")
}
