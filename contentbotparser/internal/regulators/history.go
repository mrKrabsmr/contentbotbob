package regulators

import (
	"github.com/mrkrabsmr/contentbotparser/internal/parsers"
	apiRequests "github.com/mrkrabsmr/contentbotparser/internal/requests/create"
	"log"
	"time"
)

func (r *Regulator) HistoryRegulator() {
	ch := make(chan []parsers.ParsContent)

	//time.Sleep(time.Minute * 20)

	for {
		go parsers.RgRuParser(ch)

		contents := <-ch

		err := apiRequests.NewsAPIRequest(r.Config.ApiServer, contents)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Hour)
	}

}
