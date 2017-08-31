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
	e1 := errcat.Errorf(ErrAsdf, "asdf: %s", "fmtme")
	bytes, err := json.Marshal(e1)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("must match fixture", func(t *testing.T) {
		if string(bytes) != `{"category":"err-asdf","message":"asdf: fmtme"}` {
			t.Errorf("must match fixture -- got `%s`", string(bytes))
		}
	})
	t.Run("must roundtrip", func(t *testing.T) {
		// Deserializing is interesting because if you want the category to be comparable with typeinfo,
		// you have to declare your own struct with that info.
		// Or, use a filter func to coerce it.
		type deserErr struct {
			Category_ ErrorCategory     `json:"category"          refmt:"category"`
			Message_  string            `json:"message"           refmt:"message"`
			Details_  map[string]string `json:"details,omitempty" refmt:"category,omitempty"`
		}
		var e2 deserErr
		err := json.Unmarshal(bytes, &e2)
		if err != nil {
			t.Fatal(err)
		}
		if e2.Category_ != errcat.Category(e1) {
			t.Errorf("category must match after roundtrip json -- got `%s`", e2.Category_)
		}
	})
}
