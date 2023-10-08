package parsers

import (
	"github.com/gocolly/colly"
	validators "github.com/mrkrabsmr/contentbotparser/internal"
	"strings"
	"time"
)

func RussianRTCOMParser(ch chan<- []ParsContent) {
	var contents []ParsContent

	url := "https://russian.rt.com"

	c := colly.NewCollector()

	var links []string

	c.OnHTML("li.listing__column_all-new:nth-child(n) > div:nth-child(n) > div:nth-child(n)", func(e *colly.HTMLElement) {

		str := e.DOM.Find("time").Text()
		s := strings.TrimSpace(strings.Split(str, ", ")[1]) + ":00"

		timeNow := time.Now()
		timePars, _ := time.Parse("15:04:05", s)

		mskLocation, _ := time.LoadLocation("Europe/Moscow")
		timePars = time.Date(
			timeNow.Year(),
			timeNow.Month(),
			timeNow.Day(),
			timePars.Hour(),
			timePars.Minute(),
			timePars.Second(),
			timePars.Nanosecond(),
			mskLocation,
		)

		diff := timeNow.Sub(timePars)
		if diff < time.Hour {
			link, _ := e.DOM.Find("a").Attr("href")
			links = append(links, link)
		}

	})

	c.OnScraped(func(_ *colly.Response) {
		detailC := colly.NewCollector()

		detailC.OnHTML(".article", func(e *colly.HTMLElement) {
			img, _ := e.DOM.Find("img").Attr("src")
			text := e.DOM.Find(".article__text").Text()

			if validators.TextValidator(text) {
				contentObj := ParsContent{
					Subject: "sport",
					Source:  "russianrt.com",
					Img:     img,
					Text:    text,
				}

				contents = append(contents, contentObj)
			}

		})

		for _, link := range links {
			detailC.Visit(url + link)
		}

	})

	c.Visit(url + "/tag/sport")

	ch <- contents
}
