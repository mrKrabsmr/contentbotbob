package regulators

import (
	"fmt"
	"github.com/mrkrabsmr/contentbotparser/internal/parsers"
	apiRequests "github.com/mrkrabsmr/contentbotparser/internal/requests/create"
	"log"
	"time"
)

func (r *Regulator) CryptoRegulator() {
	ch := make(chan []parsers.ParsContent)

	go func() {
		for {
			go parsers.InvestingComParser(ch)
			contents := <-ch

			err := apiRequests.NewsAPIRequest(r.Config.ApiServer, contents)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Hour)
		}
	}()

	ch2 := make(chan parsers.CryptoCompareResponse)

	go func() {
		for {
			go parsers.CryptoCompareParser(ch2)
			content := <-ch2
			fmt.Println(content)

			time.Sleep(time.Hour * 24)
		}

	}()
}
