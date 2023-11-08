package backend

const slug = "Bearer"

const (
	endRefresh            = "/refresh/"
	endRegister           = "/register/"
	endLogin              = "/login/"
	endGetProfile         = "/profile/"
	endChannelListCreate  = "/channels/"
	endChannelsUpd        = "/channels/%s/"
	endCheckChannelCreate = "/check-add/channels/"
	endSubscribes = "/subscribes/"
	endContent            = "/content/"
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
