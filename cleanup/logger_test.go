package cleanup

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func cleanup(path string) {
	os.Remove(path)
}

func TestLogTo(t *testing.T) {
	cases := []struct {
		name        string
		setup       func()
		path        string
		l           Log
		wantErr     bool
		wantNumLogs int
	}{
		{
			name:        "create new log file",
			setup:       func() {},
			path:        fakeLogFile,
			l:           fakeLog,
			wantErr:     false,
			wantNumLogs: 1,
		},
		{
			name: "append to existing log file",
			setup: func() {
				file, err := os.Create(fakeLogFile)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
				defer file.Close()
				b, err := json.Marshal(fakeLog)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
				if _, err = file.Write(b); err != nil {
					t.Fatalf("failed setup: %v", err)
				}
				if _, err = file.WriteString("\n"); err != nil {
					t.Fatalf("failed setup: %v", err)
				}
			},
			path:        fakeLogFile,
			l:           fakeLog,
			wantErr:     false,
			wantNumLogs: 2,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			t.Cleanup(func() { cleanup(tt.path) })

			err := logTo(tt.path, tt.l)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			}

			file, err := os.Open(tt.path)
			if err != nil {
				t.Fatalf("failed to verify test result: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var numLogs int
			for scanner.Scan() {
				line := scanner.Bytes()
				var validLog Log
				err := json.Unmarshal(line, &validLog)
				if err != nil {
					t.Fatalf("failed to verify test result: %v", err)
				}
				numLogs++
				continue
			}

			if numLogs != tt.wantNumLogs {
				t.Fatalf("got %d logs, want %d", numLogs, tt.wantNumLogs)
			}
		})
	}
}
