package errcat_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/polydawn/go-errcat"
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

func TestCategory(t *testing.T) {
	t.Run("uncategorized errors do not match as nil", func(t *testing.T) {
		switch errcat.Category(fmt.Errorf("womp womp")) {
		case nil:
			t.Errorf("must not nil")
		default: // pass
		}
	})
}

func TestPrefixAnnotate(t *testing.T) {
	err := errcat.ErrorDetailed(ErrAsdf, "a msg", map[string]string{"deta": "il"})
	t.Run("prefix annotation can add details", func(t *testing.T) {
		err := errcat.PrefixAnnotate(err, "more msg", [][2]string{{"more", "detail"}})
		if err.(errcat.Error).Details()["deta"] != "il" {
			t.Errorf("Lost earlier details")
		}
		if err.(errcat.Error).Details()["more"] != "detail" {
			t.Errorf("Failed to add details")
		}
	})
	t.Run("prefix annotation can prefix the message", func(t *testing.T) {
		t.Run("with basic strings", func(t *testing.T) {
			err := errcat.PrefixAnnotate(err, "more msg", [][2]string{{"more", "detail"}})
			if err.Error() != "more msg: a msg" {
				t.Fatalf("Failed to prefix message, got %q", err.Error())
			}
		})
		t.Run("with templates!", func(t *testing.T) {
			t.Run("basic interpolation", func(t *testing.T) {
				err := errcat.PrefixAnnotate(err, "using {{.tmpl}}", [][2]string{{"tmpl", "templated details"}})
				if err.Error() != "using templated details: a msg" {
					t.Errorf("Failed to prefix message, got %q", err.Error())
				}
			})
			t.Run("functions", func(t *testing.T) {
				err := errcat.PrefixAnnotate(err, "using {{.tmpl|quote}}", [][2]string{{"tmpl", "templated details"}})
				if err.Error() != "using \"templated details\": a msg" {
					t.Errorf("Failed to prefix message, got %q", err.Error())
				}
			})
		})
		t.Run("templates error referencing undefined details", func(t *testing.T) {
			err := errcat.PrefixAnnotate(err, "using {{.undefined}}", [][2]string{{"tmpl", "templated details"}})
			if err.Error() != "using <no value>: a msg" {
				t.Errorf("reference to undefined details should fail, got %q", err.Error())
			}
		})
		t.Run("templates cannot reference earlier details", func(t *testing.T) {
			err := errcat.PrefixAnnotate(err, "using {{.deta}}", [][2]string{{"tmpl", "templated details"}})
			if err.Error() != "using <no value>: a msg" {
				t.Errorf("reference to earlier details should fail, got %q", err.Error())
			}
		})
		t.Run("templates error gracefully if using a function", func(t *testing.T) {
			err := errcat.PrefixAnnotate(err, "using {{func}}", [][2]string{{"tmpl", "templated details"}})
			if err.Error() != "[[template: :1: function \"func\" not defined]]: a msg" {
				t.Errorf("reference to undefined details should fail, got %q", err.Error())
			}
		})
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
			Details_  map[string]string `json:"details,omitempty" refmt:"details,omitempty"`
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
