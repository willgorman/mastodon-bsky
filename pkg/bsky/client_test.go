package bsky_test

import (
	"context"
	"testing"

	"github.com/bluesky-social/indigo/api/atproto"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/sanity-io/litter"
	"github.com/sethvargo/go-envconfig"
	"github.com/willgorman/mastodon-bsky/pkg/bsky"
	"gotest.tools/assert"
)

func TestClientPost(t *testing.T) {
	var cfg bsky.Config
	err := envconfig.Process(context.Background(), &cfg)
	assert.NilError(t, err)

	c, err := bsky.NewClient(cfg)
	assert.NilError(t, err)

	r, err := c.Post(context.Background(), bsky.Post{
		FeedPost: appbsky.FeedPost{
			Text: "first",
			Tags: []string{"tag1", "tag2"},
			Labels: &appbsky.FeedPost_Labels{
				LabelDefs_SelfLabels: &atproto.LabelDefs_SelfLabels{
					Values: []*atproto.LabelDefs_SelfLabel{
						{
							Val: "testlabel",
						},
					},
				},
			},
		},
	})
	assert.NilError(t, err)
	t.Log(r)
}

func TestList(t *testing.T) {
	var cfg bsky.Config
	err := envconfig.Process(context.Background(), &cfg)
	assert.NilError(t, err)

	c, err := bsky.NewClient(cfg)
	assert.NilError(t, err)
	litter.Dump(c.ListRecords(context.Background()))
}
