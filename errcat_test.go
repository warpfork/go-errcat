package errcat_test

import (
	"testing"

	"."
)

func TestErrorf(t *testing.T) {
	t.Run("using string category", func(t *testing.T) {
		errcat.Errorf("catstr", "asdf: %s", "fmtme")
	})
	t.Run("using typedef string category", func(t *testing.T) {
		type ErrorCategory string
		const (
			ErrAsdf = ErrorCategory("catstrkind")
		)
		err := errcat.Errorf(ErrAsdf, "asdf: %s", "fmtme")
		switch errcat.Category(err) {
		case ErrAsdf:
			// pass
		default:
			t.Errorf("must switch")
		}
		if errcat.Category(err) != ErrAsdf {
			t.Errorf("must equal")
		}
	})
}
