package gosetta

import (
	"errors"
	"os"
)

// Rose holds translation information.
type Rose struct {
	apiKey string
}

// New creates a rose for translating.
func New() (*Rose, error) {
	apiKey := os.Getenv("TRANSLATE_API_KEY")
	if len(apiKey) < 1 {
		return nil, errors.New("missing TRANSLATE_API_KEY")
	}
	return &Rose{apiKey}, nil
}
