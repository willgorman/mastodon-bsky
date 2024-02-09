package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
)

type Config struct {
	Server       string `env:"MASTODON_SERVER"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	Username     string `env:"MASTODON_USER"`
	Password     string `env:"MASTODON_PASSWORD"`
}

type Client struct {
	*mastodon.Client
	user *mastodon.Account
}

type Status mastodon.Status

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       cfg.Server,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})

	err := c.Authenticate(ctx, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	me, err := c.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}

	return &Client{Client: c, user: me}, nil
}

func (c *Client) GetPosts(ctx context.Context) ([]*mastodon.Status, error) {
	posts, err := c.GetAccountStatuses(ctx, c.user.ID, &mastodon.Pagination{
		SinceID: mastodon.ID(""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	return posts, nil
}
