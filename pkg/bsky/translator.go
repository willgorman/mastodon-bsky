package bsky

import (
	"bytes"
	"errors"
	"strings"
	"time"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/mattn/go-mastodon"
	"golang.org/x/net/html"
	"mvdan.cc/xurls/v2"
)

const (
	RichtextFacet_Link = "app.bsky.richtext.facet#link"
	FeedPost           = "app.bsky.feed.post"
)

// translation

// TODO: (willgorman) need some config about what not to convert
// mastodon replies for example

func Convert(toot *mastodon.Status) (*appbsky.FeedPost, error) {
	// TODO: (willgorman) in order to take a mastodon image to an embedded bluesky
	// image we have to first upload the image data and get back a ref link
	// to include in the FeedPost: https://atproto.com/blog/create-post#images-embeds
	if len(toot.MediaAttachments) > 0 {
		return nil, errors.New("images not handled yet")
	}

	tootText := textContent(toot.Content)

	// TODO: (willgorman) what about Status > 300 chars?  Split into multiple FeedPost?
	// Skip it? Truncate?  Make a link back to the original mastodon post?
	return &appbsky.FeedPost{
		CreatedAt: toot.CreatedAt.Format(time.RFC3339),
		Facets:    getLinkFacets(tootText),
		Text:      tootText,
	}, nil
}

func textContent(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return s
	}
	var buf bytes.Buffer

	var extractText func(node *html.Node, w *bytes.Buffer)
	extractText = func(node *html.Node, w *bytes.Buffer) {
		if node.Type == html.TextNode {
			data := strings.Trim(node.Data, "\r\n")
			if data != "" {
				w.WriteString(data)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extractText(c, w)
		}
		if node.Type == html.ElementNode {
			name := strings.ToLower(node.Data)
			if name == "br" {
				w.WriteString("\n")
			}
		}
	}
	extractText(doc, &buf)
	return buf.String()
}

func getLinkFacets(tootText string) []*appbsky.RichtextFacet {
	urls := xurls.Relaxed()
	links := urls.FindAllStringIndex(tootText, -1)
	var facets []*appbsky.RichtextFacet
	for _, link := range links {
		facets = append(facets, &appbsky.RichtextFacet{
			Index: &appbsky.RichtextFacet_ByteSlice{
				ByteEnd:   int64(link[1]),
				ByteStart: int64(link[0]),
			},
			Features: []*appbsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Link: &appbsky.RichtextFacet_Link{
						Uri: tootText[link[0]:link[1]],
					},
				},
			},
		})
	}
	return facets
}
