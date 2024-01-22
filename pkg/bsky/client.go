package bsky

import (
	"context"
	"fmt"
	"time"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/util/cliutil"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	rpcClient *xrpc.Client
	session   *comatproto.ServerCreateSession_Output
}

type Config struct {
	PDSUrl   string `env:"PDS_URL"`
	Username string `env:"BSKY_USERNAME"`
	Password string `env:"BSKY_PASSWORD"`
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.PDSUrl == "" {
		cfg.PDSUrl = "https://bsky.social"
	}
	rpc := &xrpc.Client{
		Client: cliutil.NewHttpClient(),
		Host:   cfg.PDSUrl,
	}

	ses, err := comatproto.ServerCreateSession(context.TODO(), rpc, &comatproto.ServerCreateSession_Input{
		Identifier: cfg.Username,
		Password:   cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %w", err)
	}

	rpc.Auth = &xrpc.AuthInfo{
		AccessJwt:  ses.AccessJwt,
		RefreshJwt: ses.RefreshJwt,
		Handle:     cfg.Username,
		Did:        ses.Did,
	}

	return &Client{
		rpcClient: rpc,
		session:   ses,
	}, nil
}

type Post struct {
	appbsky.FeedPost
}

type PostResult struct {
	Cid string
	Uri string
}

func (c *Client) Post(ctx context.Context, post Post) (*PostResult, error) {
	if post.CreatedAt == "" {
		post.CreatedAt = time.Now().UTC().Format(util.ISO8601)
	}

	resp, err := comatproto.RepoCreateRecord(ctx, c.rpcClient, &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       c.rpcClient.Auth.Did,
		Record:     &lexutil.LexiconTypeDecoder{Val: &post.FeedPost},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to post: %w", err)
	}
	return &PostResult{Cid: resp.Cid, Uri: resp.Uri}, nil
}

func (c *Client) ListRecords(ctx context.Context, repoName string) ([]*bsky.FeedPost, error) {
	// TODO: (willgorman) return a channel of FeedPost backed by consuming the feed in pages
	out, err := comatproto.RepoListRecords(context.Background(),
		c.rpcClient, "app.bsky.feed.post", "", 1, repoName, true, "", "")
	if err != nil {
		return nil, fmt.Errorf("listing records: %w", err)
	}
	var ret []*bsky.FeedPost
	for _, o := range out.Records {
		rec := o.Value.Val.(*bsky.FeedPost)
		ret = append(ret, rec)
	}
	return ret, nil
}
