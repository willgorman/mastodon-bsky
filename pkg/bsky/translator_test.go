package bsky

import (
	"io"
	"net/http"
	"testing"
	"time"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/mattn/go-mastodon"
	"github.com/sanity-io/litter"
	"gotest.tools/assert"
)

var exampleMention = &mastodon.Status{
	ID:  "4321",
	URI: "https://example.com/users/me/statuses/4321",
	URL: "https://example.com/@me/4321",
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
		Emojis:         []mastodon.Emoji{}, // p0
		Moved:          nil,
		Fields:         []mastodon.Field{}, // p1
		Bot:            false,
		Discoverable:   true,
		Source:         nil,
	},
	InReplyToID:        nil,
	InReplyToAccountID: nil,
	Reblog:             nil,
	Content:            "<p>post with mention: <span class=\"h-card\" translate=\"no\"><a href=\"https://example.com/@someone\" class=\"u-url mention\">@<span>someone</span></a></span></p>",
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
	MediaAttachments:   []mastodon.Attachment{}, // p2
	Mentions: []mastodon.Mention{
		{
			URL:      "https://example.com/@someone",
			Username: "someone",
			Acct:     "someone@example.com",
			ID:       "1234",
		},
	},
	Tags: []mastodon.Tag{}, // p3
	Card: nil,
	Poll: nil,
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

var exampleImage = &mastodon.Status{
	ID:  "111795668738133188",
	URI: "https://example.com/users/me/statuses/111795668738133188",
	URL: "https://example.com/@me/111795668738133188",
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
	Content:            "<p>post with image</p>",
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
	MediaAttachments: []mastodon.Attachment{
		{
			ID:          "111795668408655718",
			Type:        "image",
			URL:         "https://files.example.com/media_attachments/files/111/795/668/408/655/718/original/499ddd2cfa22f7ef.jpg",
			RemoteURL:   "",
			PreviewURL:  "https://files.example.com/media_attachments/files/111/795/668/408/655/718/small/499ddd2cfa22f7ef.jpg",
			TextURL:     "",
			Description: "",
			Meta: mastodon.AttachmentMeta{
				Original: mastodon.AttachmentSize{
					Width:  3327,
					Height: 2493,
					Size:   "3327x2493",
					Aspect: 1.3345367027677497,
				},
				Small: mastodon.AttachmentSize{
					Width:  555,
					Height: 416,
					Size:   "555x416",
					Aspect: 1.3341346153846154,
				},
			},
		},
	},
	Mentions: []mastodon.Mention{}, // p4
	Tags:     []mastodon.Tag{},
	Card:     nil,
	Poll:     nil,
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

var exampleNewlines = &mastodon.Status{
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

var exampleTag = &mastodon.Status{
	ID:  "111795665359033449",
	URI: "https://example.com/users/me/statuses/111795665359033449",
	URL: "https://example.com/@me/111795665359033449",
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
	Content:            "<p>post with <a href=\"https://example.com/tags/tag\" class=\"mention hashtag\" rel=\"tag\">#<span>tag</span></a></p>",
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
	Tags: []mastodon.Tag{
		{
			Name:    "tag",
			URL:     "https://example.com/tags/tag",
			History: nil,
		},
	},
	Card: nil,
	Poll: nil,
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

// TODO: this doesn't actually have mastodon.Emoji. What produces that?
var exampleEmoji = &mastodon.Status{
	ID:  "111795664485035860",
	URI: "https://example.com/users/me/statuses/111795664485035860",
	URL: "https://example.com/@me/111795664485035860",
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
	Content:            "<p>post with emoji: 😀</p>",
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

var exampleLink = &mastodon.Status{
	ID:  "4321",
	URI: "https://example.com/users/me/statuses/4321",
	URL: "https://example.com/@me/4321",
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
	Content:            "<p>post with link: <a href=\"https://github.com/bluesky-social/atproto/blob/main/packages/api/README.md\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\"><span class=\"invisible\">https://</span><span class=\"ellipsis\">github.com/bluesky-social/atpr</span><span class=\"invisible\">oto/blob/main/packages/api/README.md</span></a></p>",
	CreatedAt:          time.Unix(1, 0),
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
	Card: &mastodon.Card{
		URL:          "https://github.com/bluesky-social/atproto/blob/main/packages/api/README.md",
		Title:        "atproto/packages/api/README.md at main · bluesky-social/atproto",
		Description:  "Social networking technology created by Bluesky. Contribute to bluesky-social/atproto development by creating an account on GitHub.",
		Image:        "https://files.example.com/cache/preview_cards/images/085/533/565/original/48d52fc9e782b425.jpeg",
		Type:         "link",
		AuthorName:   "",
		AuthorURL:    "",
		ProviderName: "GitHub",
		ProviderURL:  "",
		HTML:         "",
		Width:        800,
		Height:       418,
	},
	Poll: nil,
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

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		toot    *mastodon.Status
		want    *appbsky.FeedPost
		wantErr bool
	}{
		{
			name: "Toot with link",
			toot: exampleLink,
			want: &appbsky.FeedPost{
				CreatedAt: time.Unix(1, 0).Format(time.RFC3339),
				Text:      "post with link: https://github.com/bluesky-social/atproto/blob/main/packages/api/README.md",
				Facets: []*appbsky.RichtextFacet{
					{
						Features: []*appbsky.RichtextFacet_Features_Elem{
							{RichtextFacet_Link: &appbsky.RichtextFacet_Link{
								Uri: "https://github.com/bluesky-social/atproto/blob/main/packages/api/README.md",
							}},
						},
						Index: &appbsky.RichtextFacet_ByteSlice{
							ByteEnd:   90,
							ByteStart: 16,
						},
					},
				},
			},
		},
		{
			name: "Toot with tag",
			toot: exampleTag,
			want: &appbsky.FeedPost{
				CreatedAt: time.Time{}.Format(time.RFC3339),
				Text:      "post with #tag",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.toot)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.DeepEqual(t, got, tt.want)
		})
	}
}

func TestGetImg(t *testing.T) {
	resp, err := http.Get("https://files.mastodon.social/cache/preview_cards/images/085/533/565/original/48d52fc9e782b425.jpeg")
	assert.NilError(t, err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	assert.NilError(t, err)
	litter.Dump(data)
}
