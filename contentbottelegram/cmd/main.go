package main

import (
	bot "github.com/mrKrabsmr/contentbottelegram/internal"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	var config *configs.Config

	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		panic(err)
	}

	b, err := bot.NewBot(config)
	if err != nil {
		panic(err)
	}

	if err = b.Run(); err != nil {
		panic(err)
	}
}
