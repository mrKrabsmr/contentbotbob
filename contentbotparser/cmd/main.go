package main

import (
	"github.com/joho/godotenv"
	"github.com/mrkrabsmr/contentbotparser/configs"
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

	config := configs.Config{
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
