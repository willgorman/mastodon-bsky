package sync

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/willgorman/mastodon-bsky/pkg/bsky"
	"github.com/willgorman/mastodon-bsky/pkg/mastodon"
)

type mastodonSource interface {
	Open(ctx context.Context) (<-chan mastodon.Status, <-chan error)
}

type bskySink interface {
	Post(ctx context.Context, post bsky.Post) (*bsky.PostResult, error)
}

type transform func(toot *mastodon.Status) (*bsky.Post, error)

type processor struct {
	data      *Datastore
	source    mastodonSource
	sink      bskySink
	transform transform
}

func New(data *Datastore, source mastodonSource, sink bskySink) *processor {
	return &processor{
		data:   data,
		source: source,
		sink:   sink,
	}
}

func (p *processor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	toots, errors := p.source.Open(ctx)
	defer cancel()
	for {
		select {
		case toot := <-toots:
			// TODO: (willgorman)
			// add to database
			record := SyncRecord{
				AddedAt:       time.Now(),
				SourcePostID:  string(toot.ID),
				SourcePostURL: toot.URI,
			}
			if err := p.data.CreateRecord(ctx, record); err != nil {
				return fmt.Errorf("could not create sync record: %w", err)
			}
			log.Println(toot.Content)
			// convert
			post, err := p.transform(&toot)
			if err != nil {
				// TODO: (willgorman) error handling to retry on http.Get errors?
				err = fmt.Errorf("could not convert: %w", err)
				record.LastError = err.Error()
				return err
			}

			// send to sink
			result, err := p.sink.Post(ctx, *post)
			if err != nil {
				// TODO: (willgorman) retries but don't spam
				return fmt.Errorf("posting to bluesky: %w", err)
			}

			record.SyncedAt = sql.NullTime{Time: time.Now(), Valid: true}
			record.TargetPostID = result.Cid
			record.TargetPostURL = result.Uri
			err = p.data.UpdateRecord(ctx, record)
			if err != nil {
				return fmt.Errorf("failed to update after sync: %w", err)
			}

		case err := <-errors:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
