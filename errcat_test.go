package errcat_test

import (
	"encoding/json"
	"testing"

	"."
)

type ErrorCategory string

const (
	ErrAsdf = ErrorCategory("err-asdf")
)

func TestErrorf(t *testing.T) {
	t.Run("using string category", func(t *testing.T) {
		errcat.Errorf("catstr", "asdf: %s", "fmtme")
	})
	t.Run("using typedef string category", func(t *testing.T) {
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

func TestSerialization(t *testing.T) {
	err := errcat.Errorf(ErrAsdf, "asdf: %s", "fmtme")
	bytes, err := json.Marshal(err)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != `{"category":"err-asdf","message":"asdf: fmtme"}` {
		t.Errorf("must match fixture -- got `%s`", string(bytes))
	}
}
