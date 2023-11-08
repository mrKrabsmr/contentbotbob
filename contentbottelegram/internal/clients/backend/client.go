package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type APIClient struct {
	config *configs.APIClientConfig
	client *http.Client
	apiErr error
	logger *logrus.Logger
}

func NewAPIClient(config *configs.APIClientConfig, logger *logrus.Logger) *APIClient {
	return &APIClient{
		config: config,
		client: &http.Client{},
		apiErr: fmt.Errorf(""),
		logger: logger,
	}
}

func (c *APIClient) formedAuthToken(accessToken string) string {
	return slug + " " + accessToken
}

func (c *APIClient) UpdateAccessToken(refreshToken string) (*UpdateAccessTokenResponse, error) {
	url := fmt.Sprint(c.config.Address, endRefresh)
	request := &UpdateAccessTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api backend error %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	r := &UpdateAccessTokenResponse{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	err = json.Unmarshal(data, r)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return r, nil
}

func (c *APIClient) Register(request *RegisterRequest) (*RegisterResponse, error) {
	url := fmt.Sprint(c.config.Address, endRegister)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, c.apiErr
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if resp.StatusCode == http.StatusBadRequest {
		var errData map[string][]string
		var errAll string

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Error(err)
			return nil, c.apiErr
		}

		if err = json.Unmarshal(data, &errData); err != nil {
			c.logger.Error(err)
			return nil, c.apiErr
		}

		for errType := range errData {
			for _, errTxt := range errData[errType] {
				errAll += errTxt + "\n"
			}
		}

		return nil, fmt.Errorf(errAll)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, c.apiErr
	}

	return nil, nil
}

func (c *APIClient) Login(request *LoginRequest) (*LoginResponse, error) {
	url := fmt.Sprint(c.config.Address, endLogin)
	jsonData, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if resp.StatusCode == http.StatusNotFound {
		var errData map[string]string

		if err = json.Unmarshal(data, &errData); err != nil {
			c.logger.Error(err)
			return nil, c.apiErr
		}

		return nil, fmt.Errorf(errData["message"])
	}

	if resp.StatusCode == http.StatusBadRequest {
		var errData map[string][]string
		var errAll string

		if err = json.Unmarshal(data, &errData); err != nil {
			c.logger.Error(err)
			return nil, c.apiErr
		}

		for errType := range errData {
			for _, errTxt := range errData[errType] {
				errAll += errTxt + "\n"
			}
		}

		return nil, fmt.Errorf(errAll)

	}

	r := &LoginResponse{}

	if err = json.Unmarshal(data, r); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return r, nil
}

func (c *APIClient) Profile(accessToken string) (*ProfileResponse, error) {
	url := fmt.Sprint(c.config.Address, endGetProfile)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	p := &ProfileResponse{}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if err = json.Unmarshal(data, &p); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return p, nil
}

func (c *APIClient) ChannelList(accessToken string) ([]*ChannelResponse, error) {
	url := fmt.Sprint(c.config.Address, endChannelListCreate)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	var channels []*ChannelResponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	defer resp.Body.Close()

	if err = json.Unmarshal(data, &channels); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return channels, nil
}

func (c *APIClient) ChannelChangeStatus(accessToken string, channelId string, newStatus bool) error {
	url := fmt.Sprint(c.config.Address, fmt.Sprintf(endChannelsUpd, channelId))

	reqBody := &ChannelChangeStatusRequest{
		StatusOn: newStatus,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(err)
		return c.apiErr
	}

	return nil
}

func (c *APIClient) CheckChannelCreate(accessToken string) (bool, error) {
	url := fmt.Sprint(c.config.Address, endCheckChannelCreate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error(err)
		return false, c.apiErr
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return false, c.apiErr
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func (c *APIClient) ChannelCreate(accessToken string, request *ChannelCreateRequest) error {
	url := fmt.Sprint(c.config.Address, endChannelListCreate)
	jsonData, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return c.apiErr
	}

	if resp.StatusCode != http.StatusCreated {
		return c.apiErr
	}

	return nil

}

func (c *APIClient) ChannelSettings(accessToken string, request *ChannelSettingsRequest) {

}

func (c *APIClient) SubscribeList(access string) {

}

func (c *APIClient) SubscribeActivate(access string, subID string) {

}