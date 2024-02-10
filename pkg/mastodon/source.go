package mastodon

import (
	"context"
	"time"

	"github.com/mattn/go-mastodon"
)

type source struct {
	client   Client
	filters  any
	startID  string
	interval time.Duration
	statusCh chan Status
	errorCh  chan error
}

// TODO: (willgorman) define filters
func NewSource(client Client, startID string, interval time.Duration, filters any) *source {
	return &source{
		client:  client,
		filters: filters,
	}
}

func (s *source) Open(ctx context.Context) (<-chan Status, <-chan error) {
	// FIXME: (willgorman) impl
	s.statusCh = make(chan Status)
	s.errorCh = make(chan error)
	go s.streamStatus(ctx)

	return s.statusCh, s.errorCh
}

func (s *source) streamStatus(ctx context.Context) {
	var statuses []*mastodon.Status
	var err error
	tick := time.NewTicker(s.interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			close(s.errorCh)
			close(s.statusCh)
			return
		case <-tick.C:
			if len(statuses) == 0 {
				statuses, err = s.client.GetAccountStatuses(context.Background(), s.client.user.ID, &mastodon.Pagination{
					SinceID: mastodon.ID(s.startID),
				})
				if err != nil {
					s.errorCh <- err
				}
			}
			for _, status := range statuses {
				select {
				case <-ctx.Done():
					close(s.statusCh)
					close(s.errorCh)
					return
				case s.statusCh <- Status(*status):
				}
			}
			s.startID = string(statuses[len(statuses)-1].ID)
		}
	}
}

type fakeSource struct {
	toots chan Status
	errs  chan error
}

func NewFakeSource() *fakeSource {
	return &fakeSource{}
}

func (f *fakeSource) makeStatus() Status {
	exampleNewlines := &Status{
		ID:  "111795667004443647",
		URI: "https://example.com/users/me/statuses/111795667004443647",
		URL: "https://example.com/@me/111795667004443647",
		Account: mastodon.Account{
			ID:             "4321",
			Username:       "me",
			Acct:           "me",
			DisplayName:    "",
			Locked:         false,
			CreatedAt:      time.Time{},
			FollowersCount: 0,
			FollowingCount: 1,
			StatusesCount:  6,
			Note:           "<p>I use this account just for developing/testing mastodon clients</p>",
			URL:            "https://example.com/@me",
			Avatar:         "https://example.com/avatars/original/missing.png",
			AvatarStatic:   "https://example.com/avatars/original/missing.png",
			Header:         "https://example.com/headers/original/missing.png",
			HeaderStatic:   "https://example.com/headers/original/missing.png",
			Emojis:         []mastodon.Emoji{},
			Moved:          nil,
			Fields:         []mastodon.Field{},
			Bot:            false,
			Discoverable:   true,
			Source:         nil,
		},
		InReplyToID:        nil,
		InReplyToAccountID: nil,
		Reblog:             nil,
		Content:            "<p>post</p><p>with</p><p>newlines</p>",
		CreatedAt:          time.Time{},
		Emojis:             []mastodon.Emoji{},
		RepliesCount:       0,
		ReblogsCount:       0,
		FavouritesCount:    0,
		Reblogged:          false,
		Favourited:         false,
		Bookmarked:         false,
		Muted:              false,
		Sensitive:          false,
		SpoilerText:        "",
		Visibility:         "public",
		MediaAttachments:   []mastodon.Attachment{},
		Mentions:           []mastodon.Mention{},
		Tags:               []mastodon.Tag{},
		Card:               nil,
		Poll:               nil,
		Application: mastodon.Application{
			ID:           "",
			RedirectURI:  "",
			ClientID:     "",
			ClientSecret: "",
			AuthURI:      "",
		},
		Language: "en",
		Pinned:   false,
	}
	return *exampleNewlines
}

func (f *fakeSource) Open(ctx context.Context) (<-chan Status, <-chan error) {
	f.toots = make(chan Status)
	f.errs = make(chan error)
	go func() {
		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()
		for range tick.C {
			if ctx.Err() != nil {
				close(f.toots)
				close(f.errs)
				return
			}
			if f.toots == nil {
				return
			}
			time.Sleep(2 * time.Second)
			f.toots <- f.makeStatus()
		}
	}()

	return f.toots, f.errs
}
