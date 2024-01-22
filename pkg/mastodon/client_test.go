package mastodon_test

import (
	"context"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/sethvargo/go-envconfig"
	"github.com/willgorman/mastodon-bsky/pkg/mastodon"
	"gotest.tools/assert"
)

func TestGetPosts(t *testing.T) {
	var cfg mastodon.Config
	err := envconfig.Process(context.Background(), &cfg)
	assert.NilError(t, err)
	litter.Dump(cfg)

	c, err := mastodon.NewClient(context.Background(), cfg)
	assert.NilError(t, err)
	litter.Dump(c.GetPosts(context.Background()))
}
