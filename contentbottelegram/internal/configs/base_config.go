package configs

type Config struct {
	BotToken      string `yaml:"bot_token"`
    ProviderToken string `yaml:"provider_token"`
	Address       string `yaml:"address"`
	LogLevel      string `yaml:"log_level"`
	Debug         bool   `yaml:"debug"`

	SecretNum int `yaml:"secret_num"`

	Backend *APIClientConfig `yaml:"backend"`
	Redis   *RedisConfig     `yaml:"redis"`
}
