package backend

const (
	endRefresh           = "/refresh/"
	endRegister          = "/register/"
	endLogin             = "/login/"
	endGetProfile        = "/profile/"
	endChannelListCreate = "/channels/"
	endChannelsUpd       = "/channels/<slug:id>/"
	endContent           = "/content/"
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
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	RegisteredAt string      `json:"registered_at"`
	Subscribe    interface{} `json:"subscribe"`
}

type ChannelListResponse struct {
}

type ChannelCreateRequest struct {
}

type ChannelSettingsRequest struct {
}
