package sync

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const create string = `
	CREATE TABLE IF NOT EXISTS sync_record (
		id INTEGER NOT NULL PRIMARY KEY,
		added_at DATETIME NOT NULL,
		synced_at DATETIME NULL,
		source_post_id TEXT,
		source_post_url TEXT,

		target_post_id TEXT,
		target_post_url TEXT
	);
`

type SyncRecord struct {
	ID            int          `db:"id"`
	AddedAt       time.Time    `db:"added_at"`
	SyncedAt      sql.NullTime `db:"synced_at"`
	SourcePostID  string       `db:"source_post_id"`
	SourcePostURL string       `db:"source_post_url"`
	TargetPostID  string       `db:"target_post_id"`
	TargetPostURL string       `db:"target_post_url"`
}

func CreateDatastore(path string) error {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	_, err = db.Exec(create)
	if err != nil {
		return fmt.Errorf("failed to exec create: %w", err)
	}

	return nil
}

type Datastore struct {
	db *sqlx.DB
}

func OpenDatastore(path string) (*Datastore, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	return &Datastore{db: db}, nil
}

func (d *Datastore) ListRecords(ctx context.Context) ([]SyncRecord, error) {
	var records []SyncRecord
	err := d.db.SelectContext(ctx, &records, `SELECT * FROM sync_record`)
	if err != nil {
		return nil, fmt.Errorf("unable to query records: %w", err)
	}
	return records, nil
}

func (d *Datastore) GetRecord(ctx context.Context) (*SyncRecord, error) {
	return nil, nil
}

func (d *Datastore) CreateRecord(ctx context.Context, record SyncRecord) error {
	if record.AddedAt.IsZero() {
		record.AddedAt = time.Now().UTC()
	}
	_, err := d.db.NamedExecContext(ctx,
		`INSERT INTO sync_record (added_at, synced_at, source_post_id, source_post_url, target_post_id, target_post_url)
			VALUES (:added_at, :synced_at, :source_post_id, :source_post_url, :target_post_id, :target_post_url)
		`, &record)
	return err
}

func (d *Datastore) UpdateRecord(ctx context.Context, record SyncRecord) error {
	return nil
}
