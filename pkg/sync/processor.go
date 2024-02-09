package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/willgorman/mastodon-bsky/pkg/mastodon"
)

type mastodonSource interface {
	Open() (<-chan mastodon.Status, <-chan error)
	Close() error
}

type bskySink interface{}

type processor struct {
	data   *Datastore
	source mastodonSource
	sink   bskySink
}

func New(data *Datastore, source mastodonSource, sink bskySink) *processor {
	return &processor{
		data:   data,
		source: source,
		sink:   sink,
	}
}

func (p *processor) Run(ctx context.Context) error {
	toots, errors := p.source.Open()
	defer p.source.Close()
	for {
		select {
		case toot := <-toots:
			// TODO: (willgorman)
			// add to database
			if err := p.data.CreateRecord(ctx, SyncRecord{
				AddedAt:       time.Now(),
				SourcePostID:  string(toot.ID),
				SourcePostURL: toot.URI,
			}); err != nil {
				return fmt.Errorf("could not create sync record: %w", err)
			}

			// convert
			// send to sink

			log.Println(toot.Content)
		case err := <-errors:
			// TODO: (willgorman) shutdown stuff
			return err
		case <-ctx.Done():
			// TODO: (willgorman) shutdown stuff
			return ctx.Err()
		}
	}
}
