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

func (c *APIClient) ChannelList() ([]*ChannelResponse, error) {
	url := fmt.Sprint(c.config.Address, fmt.Sprintf(endChannels, 1))

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	var channels []*ChannelResponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Ошибка от сервера, код: ", resp.StatusCode)
		return nil, c.apiErr
	}

	if err = json.Unmarshal(data, &channels); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return channels, nil
}

func (c *APIClient) ChannelsUser(accessToken string) ([]*ChannelResponse, error) {
	url := fmt.Sprint(c.config.Address, endChannelListCreate)

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

	var channels []*ChannelResponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("необходимо подтвердить аккаунт в профиле")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Ошибка от сервера, код: ", resp.StatusCode)
		return nil, c.apiErr
	}

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

func (c *APIClient) ChannelSettings(accessToken string, chanId string) (*ChannelSettingsResponse, error) {
	url := fmt.Sprint(c.config.Address, fmt.Sprintf(endChannelSettingsUpd, chanId))

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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Ошибка от сервера, код: ", resp.StatusCode)
		return nil, c.apiErr
	}

	settings := &ChannelSettingsResponse{}

	if err = json.Unmarshal(data, settings); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return settings, nil

}

func (c *APIClient) ChannelSettingsUpdate(accessToken string, request *ChannelSettingsChangeRequest) error {
	url := fmt.Sprint(c.config.Address, fmt.Sprintf(endChannelSettingsUpd, request.ID))

	jsonData, err := json.Marshal(request.Data)
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
		c.logger.Errorf("Ошибка от сервера (обновление настроек): %d", resp.StatusCode)
		return c.apiErr
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(err)
		return c.apiErr
	}

	return nil
}

func (c *APIClient) SubscribeList(access string) (*SubscribesResponse, error) {
	url := fmt.Sprint(c.config.Address, endSubscribes)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}
	req.Header.Set("Authorization", c.formedAuthToken(access))

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("Ошибка от сервера (получение списка подписок): %d", resp.StatusCode)
		return nil, c.apiErr
	}

	subscribes := &SubscribesResponse{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	if err = json.Unmarshal(data, subscribes); err != nil {
		c.logger.Error(err)
		return nil, c.apiErr
	}

	return subscribes, nil
}

func (c *APIClient) SubscribeActivate(access string, subID string) {

}

func (c *APIClient) Content(chanId string) (*ContentInfoResponse, bool, error) {
	url := fmt.Sprint(c.config.Address, fmt.Sprintf(endContent, chanId, c.config.APIKey))

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, false, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	content := &ContentInfoResponse{}

	if err = json.Unmarshal(data, content); err != nil {
		return nil, false, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("ошибка от сервера: %v", resp.StatusCode)
	}

	if !content.Active {
		return nil, false, nil
	}

	return content, true, nil
}

func (c *APIClient) SendActivateCode(accessToken string) error {
	url := fmt.Sprint(c.config.Address, endSendActivateCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.formedAuthToken(accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("Ошибка от сервера (отправка кода на почту): %d", resp.StatusCode)
		return c.apiErr
	}

	return nil
}

func (c *APIClient) UserEmailActivate(accessToken string, request *UserEmailActivateRequest) error {
	url := fmt.Sprint(c.config.Address, endUserActivate)
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

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("перед подтверждением необходимо сначала получить сообщение на почту")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("код введен неверно")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(err)
		return c.apiErr
	}

	return nil
}
