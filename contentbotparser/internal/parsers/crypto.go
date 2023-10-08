package parsers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	validators "github.com/mrkrabsmr/contentbotparser/internal"
	"net/http"
	"strings"
)

func InvestingComParser(ch chan<- []ParsContent) {
	var contents []ParsContent

	url := "https://ru.investing.com/news/cryptocurrency-news"
	c := colly.NewCollector()

	var links []string
	c.OnHTML(".largeTitle > article:nth-child(n) > div:nth-child(n)", func(e *colly.HTMLElement) {
		e.ForEach("*", func(i int, element *colly.HTMLElement) {

			if element.Name == "span" {
				element.ForEach("span", func(o int, el *colly.HTMLElement) {
					if strings.Contains(el.Text, "-") {
						parts := strings.Split(el.Text[5:], " ")

						if len(parts) > 1 {
							time := parts[1]

							if strings.Contains(time, "минут") || strings.Contains(time, "секунд") {
								if element.DOM.Parent().Find("a").Length() > 0 {
									link, _ := element.DOM.Parent().Find("a").Attr("href")
									links = append(links, url+link)
								}
							}
						}
					}

				})
			}
		})

	})

	c.OnScraped(func(_ *colly.Response) {
		detailC := colly.NewCollector()

		detailC.OnHTML("div.WYSIWYG:nth-child(8)", func(e *colly.HTMLElement) {
			txt := strings.TrimPrefix(
				strings.TrimPrefix(strings.Split(e.DOM.Find("p").Text(), "Читайте")[0], "Позиция успешно добавлена: \n"), "Happycoin.club - ",
			)

			if validators.TextValidator(txt) {
				if !strings.HasPrefix(txt, "Investing.com") {
					var contentObj ParsContent
					contentObj.Subject = "crypto"
					contentObj.Text = txt
					contentObj.Source = "investing.com"
					e.ForEach("img", func(i int, element *colly.HTMLElement) {
						contentObj.Img = element.Attr("src")
					})

					contents = append(contents, contentObj)
				}
			}
		})

		for _, link := range links {
			detailC.Visit(link)
		}
	})

	c.Visit(url)

	ch <- contents
}

func CryptoCompareParser(ch chan<- CryptoCompareResponse) {
	url := "https://min-api.cryptocompare.com/data/top/totalvolfull?limit=10&tsym=USD"
	response, _ := http.Get(url)

	var r CryptoCompareResponse
	json.NewDecoder(response.Body).Decode(&r)

	ch <- r
}
