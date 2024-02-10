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
		source_post_id TEXT PRIMARY KEY,
		source_post_url TEXT,
		target_post_id TEXT,
		target_post_url TEXT,
		added_at DATETIME NOT NULL,
		synced_at DATETIME NULL,
		attempts INT DEFAULT 0 NOT NULL,
		last_error TEXT DEFAULT "" NOT NULL
	);
`

type SyncRecord struct {
	AddedAt       time.Time    `db:"added_at"`
	SyncedAt      sql.NullTime `db:"synced_at"`
	SourcePostID  string       `db:"source_post_id"`
	SourcePostURL string       `db:"source_post_url"`
	TargetPostID  string       `db:"target_post_id"`
	TargetPostURL string       `db:"target_post_url"`
	Attempts      int          `db:"attempts"`
	LastError     string       `db:"last_error"`
}

func CreateDatastore(path string) (*Datastore, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	_, err = db.Exec(create)
	if err != nil {
		return nil, fmt.Errorf("failed to exec create: %w", err)
	}

	return &Datastore{db: db}, nil
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

func (d *Datastore) GetRecord(ctx context.Context, sourcePostID string) (*SyncRecord, error) {
	record := SyncRecord{}
	err := d.db.GetContext(ctx, &record,
		`SELECT * FROM sync_record WHERE source_post_id = ?`, sourcePostID)
	return &record, err
}

func (d *Datastore) CreateRecord(ctx context.Context, record SyncRecord) error {
	if record.AddedAt.IsZero() {
		record.AddedAt = time.Now().UTC()
	}
	_, err := d.db.NamedExecContext(ctx,
		`INSERT INTO sync_record (added_at, synced_at, source_post_id, source_post_url, target_post_id, target_post_url, attempts)
			VALUES (:added_at, :synced_at, :source_post_id, :source_post_url, :target_post_id, :target_post_url, 0)
		`, &record)
	return err
}

func (d *Datastore) UpdateRecord(ctx context.Context, record SyncRecord) error {
	_, err := d.db.NamedExecContext(ctx,
		`UPDATE sync_record 
			SET synced_at = CURRENT_TIMESTAMP, 
					target_post_id = :target_post_id, 
					target_post_url = :target_post_url,
					last_error = :last_error,
					attempts = attempts+1`, &record)
	return err
}
