package apiRequests

import (
	"bytes"
	"encoding/json"
	"github.com/mrkrabsmr/contentbotparser/internal/parsers"
	"net/http"
)

func NewsAPIRequest(apiUrl string, data []parsers.ParsContent) error {
	endpoint := "/parser/contents/news/"

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(apiUrl+endpoint, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	return nil
}
