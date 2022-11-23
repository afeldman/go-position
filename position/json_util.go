package position

import (
	"encoding/json"
	"errors"
	"io"
)

func jsonType(in io.Reader) (int, error) {
	dec := json.NewDecoder(in)
	// Get just the first valid JSON token from input
	t, err := dec.Token()
	if err != nil {
		return -1, err
	}
	if d, ok := t.(json.Delim); ok {
		// The first token is a delimiter, so this is an array or an object
		switch d {
		case '[':
			return 0, nil
		case '{':
			return 1, nil
		default:
			return -1, errors.New("Unexpected delimiter")
		}
	}
	return -1, errors.New("Input does not represent a JSON object or array")
}
