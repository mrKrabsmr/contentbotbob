package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/mrKrabsmr/contentbottelegram/internal/configs"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	config *configs.RedisConfig
	*redis.Client
}

func NewClient(config *configs.RedisConfig) *Client {
	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
		DB:   config.DB,
	}

	return &Client{
		config,
		redis.NewClient(opts),
	}
}

func (c *Client) GetTokens(userId int64) (string, string, error) {
	keyAccess := fmt.Sprintf("%v:%v", userId, "access")
	keyRefresh := fmt.Sprintf("%v:%v", userId, "refresh")

	refreshToken, errRefresh := c.Get(context.Background(), keyRefresh).Result()
	if errRefresh != nil {
		if !errors.Is(errRefresh, redis.Nil) {
			return "", "", errRefresh
		}

		return "", "", nil
	}

	accessToken, errAccess := c.Get(context.Background(), keyAccess).Result()
	if errAccess != nil {
		if !errors.Is(errAccess, redis.Nil) {
			return "", "", errAccess
		}
	}

	return accessToken, refreshToken, nil
}

func (c *Client) SetAccessToken(userId int64, accessToken string) error {
	keyAccess := fmt.Sprintf("%v:%v", userId, "access")

	lt := time.Duration(c.config.AccessTokenLT)
	_, err := c.Set(context.Background(), keyAccess, accessToken, lt).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetRefreshToken(userId int64, refreshToken string) error {
	keyRefresh := fmt.Sprintf("%v:%v", userId, "refresh")

	lt := time.Duration(c.config.RefreshTokenLT)
	_, err := c.Set(context.Background(), keyRefresh, refreshToken, lt).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteTokens(userId int64) error {
	keyAccess := fmt.Sprintf("%v:%v", userId, "access")
	keyRefresh := fmt.Sprintf("%v:%v", userId, "refresh")

	_, err := c.Del(context.Background(), keyAccess).Result()
	if err != nil {
		return err
	}

	_, err = c.Del(context.Background(), keyRefresh).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetPostFormat(chanId string, toChannel bool) error {
	_, err := c.Set(context.Background(), chanId, toChannel, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
