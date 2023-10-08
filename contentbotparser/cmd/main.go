package main

import (
	"github.com/joho/godotenv"
	config2 "github.com/mrkrabsmr/contentbotparser/config"
	"github.com/mrkrabsmr/contentbotparser/internal/regulators"
	"log"
	"os"
	"sync"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config := config2.Config{
		ApiServer: os.Getenv("API_SERVER"),
	}

	regulator := &regulators.Regulator{
		Config: &config,
	}

	wg := sync.WaitGroup{}

	wg.Add(3)
	go regulator.CryptoRegulator()
	go regulator.SportRegulator()
	go regulator.HistoryRegulator()

	wg.Wait()

}
