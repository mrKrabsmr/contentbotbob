package configs

type Config struct {
	BotToken string `yaml:"bot_token"`

	Address  string `yaml:"address"`
	LogLevel string `yaml:"log_level"`
	Debug    bool   `yaml:"debug"`

	Backend *APIClientConfig `yaml:"backend"`
	Redis   *RedisConfig     `yaml:"redis"`
}
