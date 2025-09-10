package cleanup

import (
	"encoding/json"
	"os"
)

func logTo(path string, l Log) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	b = append(b, '\n')

	if _, err = file.Write(b); err != nil {
		return err
	}

	return nil
}
