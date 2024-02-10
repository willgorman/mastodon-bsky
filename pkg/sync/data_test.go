package sync

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/sanity-io/litter"
	"gotest.tools/assert"
)

func TestCreateDatabase(t *testing.T) {
	dir := t.TempDir()
	_, err := CreateDatastore(fmt.Sprintf("%s/sync.db", dir))
	assert.NilError(t, err)
}

func TestRecords(t *testing.T) {
	dir := t.TempDir()
	path := fmt.Sprintf("%s/sync.db", dir)
	_, err := CreateDatastore(path)
	assert.NilError(t, err)

	ds, err := OpenDatastore(path)
	assert.NilError(t, err)
	err = ds.CreateRecord(context.Background(), SyncRecord{
		SourcePostID: "a",
		TargetPostID: "b",
	})
	assert.NilError(t, err)
	records, err := ds.ListRecords(context.Background())
	assert.NilError(t, err)
	litter.Dump(records)

	record, err := ds.GetRecord(context.Background(), "a")
	assert.NilError(t, err)
	litter.Dump(record)

	record.LastError = "oh no"
	err = ds.UpdateRecord(context.Background(), *record)
	assert.NilError(t, err)
	litter.Dump(ds.GetRecord(context.Background(), "a"))

	record.LastError = ""
	record.TargetPostID = "c"
	record.TargetPostURL = "http://example.com"
	err = ds.UpdateRecord(context.Background(), *record)
	assert.NilError(t, err)
	litter.Dump(ds.GetRecord(context.Background(), "a"))
}

func init() {
	timeType := reflect.TypeOf(time.Time{})
	litter.Config.DumpFunc = func(v reflect.Value, w io.Writer) bool {
		if v.Type() != timeType {
			return false
		}

		t := v.Interface().(time.Time)
		fmt.Fprintf(w, `{/* %s */}`, t.Format(time.RFC3339))
		return true
	}
}
