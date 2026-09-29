package postgres

import (
	"context"
	"errors"
	"testing"

	"hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/modules/courses/domain"
)

func TestTranslate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input error
		want  error
	}{
		{"not found", &ent.NotFoundError{}, domain.ErrNotFound},
		{"constraint", &ent.ConstraintError{}, domain.ErrConflict},
		{"validation", &ent.ValidationError{}, domain.ErrInvalid},
		{"cancellation", context.Canceled, context.Canceled},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := translate(tc.input); !errors.Is(got, tc.want) {
				t.Fatalf("translate(%v) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}
