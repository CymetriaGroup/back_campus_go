package identifier_test

import (
	"testing"

	"hexagonal-go-backend/internal/platform/identifier"
)

func TestNewReturnsValidUniqueULIDs(t *testing.T) {
	first := identifier.New()
	second := identifier.New()

	if err := identifier.Validate(first); err != nil {
		t.Fatalf("Validate(%q): %v", first, err)
	}
	if err := identifier.Validate(second); err != nil {
		t.Fatalf("Validate(%q): %v", second, err)
	}
	if first == second {
		t.Fatalf("expected unique ULIDs, got %q twice", first)
	}
}

func TestValidateRejectsNonCanonicalULIDs(t *testing.T) {
	valid := identifier.New()
	for _, value := range []string{"", "not-a-ulid", valid[:identifier.Length-1]} {
		if err := identifier.Validate(value); err == nil {
			t.Fatalf("expected %q to be rejected", value)
		}
	}
}
