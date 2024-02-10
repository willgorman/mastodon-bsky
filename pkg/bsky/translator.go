package bsky

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

type Card struct {
	appbsky.EmbedExternal_External
	ThumbImg io.ReadCloser
}
type Post struct {
	appbsky.FeedPost
	images map[*appbsky.EmbedImages_Image]io.ReadCloser
	card   Card
}

// translation

// TODO: (willgorman) need some config about what not to convert
// mastodon replies for example
// TODO: (willgorman) change return type, add image/embed data to the struct
func Convert(toot *mastodon.Status) (*Post, error) {
	result := &Post{}
	tootText := textContent(toot.Content)

	// TODO: (willgorman) what about Status > 300 chars?  Split into multiple FeedPost?
	// Skip it? Truncate?  Make a link back to the original mastodon post?
	result.FeedPost = appbsky.FeedPost{
		CreatedAt: toot.CreatedAt.Format(time.RFC3339),
		Facets:    getLinkFacets(tootText),
		Text:      tootText,
	}
	// in order to take a mastodon image to an embedded bluesky
	// image we have to first upload the image data and get back a ref link
	// to include in the FeedPost: https://atproto.com/blog/create-post#images-embeds
	// so we'll fetch the image from the source url here and pass the data along
	// for the bluesky client to upload and attach the link to the EmbedImages_Image
	for _, attachment := range toot.MediaAttachments {
		if attachment.Type != "image" {
			continue
		}
		if result.images == nil {
			result.images = make(map[*appbsky.EmbedImages_Image]io.ReadCloser)
		}
		resp, err := http.Get(attachment.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to get image from %s: %w", attachment.URL, err)
		}
		result.images[&appbsky.EmbedImages_Image{
			Alt: attachment.Description,
			AspectRatio: &appbsky.EmbedImages_AspectRatio{
				Height: attachment.Meta.Original.Height,
				Width:  attachment.Meta.Original.Width,
			},
		}] = resp.Body
	}
	for image := range result.images {
		result.Embed.EmbedImages.Images = append(result.Embed.EmbedImages.Images, image)
	}

	if toot.Card != nil {
		result.card = Card{
			EmbedExternal_External: appbsky.EmbedExternal_External{
				Description: toot.Card.Description,
				Title:       toot.Card.Title,
				Uri:         toot.Card.URL,
			},
		}
		res, err := http.Get(toot.Card.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch card image from %s: %w", toot.Card.Image, err)
		}
		result.card.ThumbImg = res.Body
		result.Embed.EmbedExternal.External = &result.card.EmbedExternal_External
	}

	return result, nil
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
