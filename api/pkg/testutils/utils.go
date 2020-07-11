package testutils

import (
	"bytes"
	"encoding/json"
)

// FormatJSON formats json string to be added to golden file
func FormatJSON(b []byte) (string, error) {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, b, "", "\t")
	if err != nil {
		return "", err
	}
	return string(formatted.Bytes()), nil
}
