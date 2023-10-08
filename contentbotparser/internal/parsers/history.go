package parsers

import (
	"github.com/gocolly/colly"
	validators "github.com/mrkrabsmr/contentbotparser/internal"
)

func RgRuParser(ch chan<- []ParsContent) {
	var contents []ParsContent

	url := "https://rg.ru/"

	c := colly.NewCollector()

	var links []string

	c.OnHTML(
		"div.PageRubricContent_listItem__rjCcF:nth-child(n) > div:nth-child(n) > div:nth-child(2) > a:nth-child(n)",
		func(e *colly.HTMLElement) {
			if len(e.Text) == 5 {
				link := e.Attr("href")
				links = append(links, link)
			}
		})

	c.OnScraped(func(_ *colly.Response) {
		detailC := colly.NewCollector()

		detailC.OnHTML(".PageArticle_main__eAtJy", func(e *colly.HTMLElement) {
			img := e.ChildAttr(".PageArticleHead_image__fR3_1", "src")
			text := e.ChildText(".PageArticleContentStyling_text__zUBlk")

			if validators.TextValidator(text) {
				contentObj := ParsContent{
					Subject: "history",
					Source:  "rg.ru",
					Img:     img,
					Text:    text,
				}

				contents = append(contents, contentObj)
			}
		})

		for _, link := range links {
			newURL := url + link
			detailC.Visit(newURL)
		}
	})

	c.Visit(url + "tema/obshestvo/istorija")

	ch <- contents
}
