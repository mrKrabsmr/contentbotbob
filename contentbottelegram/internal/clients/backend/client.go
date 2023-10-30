package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"io"
	"net/http"
)

type APIClient struct {
	config *configs.APIClientConfig
	client *http.Client
	apiErr error
}

func NewAPIClient(config *configs.APIClientConfig) *APIClient {
	return &APIClient{
		config: config,
		client: &http.Client{},
		apiErr: fmt.Errorf(""),
	}
}

func (c *APIClient) UpdateAccessToken(refreshToken string) (*UpdateAccessTokenResponse, error) {
	url := fmt.Sprint(c.config.Address, endRefresh)
	request := &UpdateAccessTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api backend error %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	r := &UpdateAccessTokenResponse{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *APIClient) Register(request *RegisterRequest) (*RegisterResponse, error) {
	url := fmt.Sprint(c.config.Address, endRegister)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		var errors map[string][]string
		var errAll string

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, c.apiErr
		}

		defer resp.Body.Close()

		if err = json.Unmarshal(data, &errors); err != nil {
			return nil, c.apiErr
		}

		for errType := range errors {
			for _, errTxt := range errors[errType] {
				errAll += errTxt + "\n"
			}
		}

		return nil, fmt.Errorf(errAll)
	}

	if resp.StatusCode != 201 {
		return nil, c.apiErr
	}

	return nil, nil
}

func (c *APIClient) Login(request *LoginRequest) (*LoginResponse, error) {
	url := fmt.Sprint(c.config.Address, endLogin)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r := &LoginResponse{}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *APIClient) Profile(accessToken string) (*ProfileResponse, error) {
	url := fmt.Sprint(c.config.Address, endGetProfile)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	p := &ProfileResponse{}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	return p, nil
}

//
//func (c *APIClient) ChannelList(accessToken string) (*ChannelListResponse, error) {
//
//}
//
//func (c *APIClient) ChannelCreate(accessToken string, request *ChannelCreateRequest) error {
//
//}
//
//func (c *APIClient) ChannelSettings(accessToken string, request *ChannelSettingsRequest) error {
//
//}
