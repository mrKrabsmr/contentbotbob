package backend

import "encoding/json"

const slug = "Bearer"

const (
	endRefresh            = "/refresh/"
	endRegister           = "/register/"
	endLogin              = "/login/"
	endGetProfile         = "/profile/"
	endChannelListCreate  = "/channels/"
	endChannels           = "/list/channels?status_on=%d"
	endChannelsUpd        = "/channels/%s/"
	endChannelSettingsUpd = "/channel-settings/%s/"
	endCheckChannelCreate = "/check-add/channels/"
	endSubscribes         = "/subscribes/"
	endContent            = "/contents/for/%s?key=%s"
	endSendActivateCode   = "/send/activate-code/"
	endUserActivate       = "/users/activate/"
)

type UpdateAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ProfileResponse struct {
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	RegisteredAt string                 `json:"registered_at"`
	IsConfirmed  bool                   `json:"is_confirmed"`
	Subscribe    map[string]interface{} `json:"subscribe"`
}

type ChannelResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Type     string `json:"type"`
	Resource string `json:"resource"`
	OuterID  string `json:"outer_id"`
	StatusOn bool   `json:"status_on"`
}

type ChannelChangeStatusRequest struct {
	StatusOn bool `json:"status_on"`
}

type ChannelCreateRequest struct {
	Name     string                  `json:"name"`
	Type     string                  `json:"type"`
	Resource string                  `json:"resource"`
	OuterID  string                  `json:"outer_id"`
	Settings *ChannelSettingsRequest `json:"settings"`
}

type ChannelSettingsRequest struct {
	MinRating        int    `json:"min_rating"`
	EmptyFileAllowed bool   `json:"empty_file_allowed"`
	NotLaterDays     int    `json:"not_later_days"`
	Language         string `json:"language"`
}

type ChannelSettingsResponse struct {
	ID               string `json:"id"`
	MinRating        int    `json:"min_rating"`
	EmptyFileAllowed bool   `json:"empty_file_allowed"`
	NotLaterDays     int    `json:"not_later_days"`
	Language         string `json:"language"`
}

type ChannelSettingsChangeRequest struct {
	ID   string
	Data map[string]interface{}
}

type ContentImagesResponse struct {
	ID      string `json:"id"`
	FileUrl string `json:"file_url"`
}

type ContentResponse struct {
	ID     string                   `json:"id"`
	Text   string                   `json:"text"`
	Images []*ContentImagesResponse `json:"images"`
}

type ContentInfoResponse struct {
	Result *ContentResponse `json:"result"`
	Active bool             `json:"active"`
}

type UserEmailActivateRequest struct {
	Code string `json:"code"`
}

type SubscribeResponse struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	MaxUseChannels int         `json:"max_use_channels"`
	PriceRub       json.Number `json:"price_rub,omitempty"`
}

type SubscribesResponse struct {
	Result []*SubscribeResponse `json:"result"`
}
