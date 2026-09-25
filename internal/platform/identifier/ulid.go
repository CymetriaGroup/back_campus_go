package identifier

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

const Length = 26

// New returns a canonical, lexicographically sortable ULID.
func New() string {
	return ulid.Make().String()
}

// Validate checks that value is a canonical ULID representation.
func Validate(value string) error {
	parsed, err := ulid.ParseStrict(value)
	if err != nil || parsed.String() != value {
		return fmt.Errorf("invalid ULID")
	}
	return nil
}
