package cleanup

import (
	"fmt"
	"log"
	"sync"
)

func FailedSave(staging MetadataRemover, mr MediaRemover, nr NotesRemover, uid string, oid string, options ...*Options) {
	var opts *Options
	if len(options) == 0 {
		opts = nil
	} else {
		opts = options[0]
	}

	targets := []struct {
		name    string
		cleanup func(uid string, oid string) error
	}{
		{
			name: "staged",
			cleanup: func(uid, oid string) error {
				return staging.MetadataRemove(uid, oid)
			},
		},
		{
			name: "media",
			cleanup: func(uid, oid string) error {
				return mr.MediaRemove(uid, oid)
			},
		},
		{
			name: "notes",
			cleanup: func(uid, oid string) error {
				return nr.NotesRemove(uid, oid)
			},
		},
	}

	var wg sync.WaitGroup

	for _, target := range targets {
		if target.name == "media" && !opts.HasMedia {
			continue
		}

		wg.Add(1)

		go func() {
			err := target.cleanup(uid, oid)
			if err != nil {
				l := Log{Uid: uid, Id: oid, Err: err.Error()}
				logPath := fmt.Sprintf("failed_%s_remove.txt", target.name)
				err := logTo(logPath, l)
				if err != nil {
					log.Printf("failed to log %q to %s: %v", logPath, l, err)
				}
				wg.Done()
			}
		}()
	}

	wg.Wait()
}
