package configs

type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	DB   int    `yaml:"db"`

	AccessTokenLT  int64 `yaml:"access_token_lifetime"`
	RefreshTokenLT int64 `yaml:"refresh_token_lifetime"`
}
