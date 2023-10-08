package apiRequests

import (
	"bytes"
	"encoding/json"
	"github.com/mrkrabsmr/contentbotparser/internal/parsers"
	"net/http"
)

func NewsAPIRequest(apiUrl string, data []parsers.ParsContent) error {
	endpoint := "parser/content/news/"

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	http.Post(apiUrl+endpoint, "application/json", bytes.NewReader(jsonData))

	return nil
}
