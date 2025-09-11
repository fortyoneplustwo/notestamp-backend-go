package cleanup

import (
	"errors"
	"io/fs"
	"notestamp/media"
	"notestamp/metadata"
	"notestamp/notes"
	"os"
	"testing"
)

var fakeError = errors.New("fakeError")

func TestFailedSave(t *testing.T) {
	cases := []struct {
		name    string
		staging metadata.MockMetadataStore
		mr      media.MockMediaStore
		nr      notes.MockNotesStore
		uid     string
		oid     string
		options *Options
	}{
		{
			name:    "failed cleanup with media",
			staging: metadata.MockMetadataStore{Err: fakeError},
			mr:      media.MockMediaStore{Err: fakeError},
			nr:      notes.MockNotesStore{Err: fakeError},
			uid:     metadata.FakeUid,
			oid:     metadata.FakeMetadata.Title,
			options: &Options{HasMedia: true},
		},
		{
			name:    "failed cleanup no media",
			staging: metadata.MockMetadataStore{Err: fakeError},
			mr:      media.MockMediaStore{Err: fakeError},
			nr:      notes.MockNotesStore{Err: fakeError},
			uid:     metadata.FakeUid,
			oid:     metadata.FakeMetadata.Title,
			options: &Options{HasMedia: true},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				if tt.options.HasMedia {
					os.Remove("failed_media_remove.txt")
				}
				os.Remove("failed_notes_remove.txt")
				os.Remove("failed_staged_remove.txt")
			})

			FailedSave(tt.staging, tt.mr, tt.nr, tt.uid, tt.oid, tt.options)

			if tt.options.HasMedia {
				_, err := os.Stat("failed_media_remove.txt")
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						t.Fatalf("media log file not created")
					} else {
						t.Fatal(err)
					}
				}
			}

			_, err := os.Stat("failed_notes_remove.txt")
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("notes log file not created")
				} else {
					t.Fatal(err)
				}
			}

			_, err = os.Stat("failed_staged_remove.txt")
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("staged log file not created")
				} else {
					t.Fatal(err)
				}
			}
		})
	}
}
