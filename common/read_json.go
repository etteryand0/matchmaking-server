package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func ReadJson(path string, v any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintln("error opening", path))
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return errors.New(fmt.Sprintln("error reading", path))
	}
	if err = json.Unmarshal(byteValue, v); err != nil {
		return errors.New(fmt.Sprintln("error unmarshalling", path))
	}

	return nil
}
