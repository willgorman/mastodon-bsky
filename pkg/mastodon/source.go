package mastodon

import (
	"time"

	"github.com/mattn/go-mastodon"
)

type source struct{}

// TODO: (willgorman) define filters
func NewSource(client Client, filters any) source {
	return source{}
}

func (s *source) Open() (<-chan Status, <-chan error) {
	// FIXME: (willgorman) impl
	return nil, nil
}

func (s *source) Close() error {
	// FIXME: (willgorman) impl
	return nil
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

func (f *fakeSource) Open() (<-chan Status, <-chan error) {
	f.toots = make(chan Status)
	f.errs = make(chan error)
	go func() {
		for {
			if f.toots == nil {
				return
			}
			time.Sleep(2 * time.Second)
			f.toots <- f.makeStatus()
		}
	}()

	return f.toots, f.errs
}

func (f *fakeSource) Close() error {
	// TODO: (willgorman) impl
	close(f.toots)
	f.toots = nil
	close(f.errs)
	f.errs = nil
	return nil
}
