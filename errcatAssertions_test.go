package errcat_test

import (
	"fmt"
	"testing"

	"github.com/polydawn/go-errcat"
)

type ErrorCategoryA string
type ErrorCategoryB string

const (
	ErrQwer = ErrorCategoryA("err-qwer")
	ErrZxcv = ErrorCategoryB("err-zxcv")
)

func TestDeferredAssertion(t *testing.T) {
	t.Run("assertions silent on nil errors", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			return nil
		}()
		shouldCategory(t, err, nil)
	})
	t.Run("using string category", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			return func() (err error) {
				return errcat.Errorf(ErrQwer, "aaah")
			}()
		}()
		shouldCategory(t, err, ErrQwer)
	})
	t.Run("using string category", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			err = func() (err error) {
				defer errcat.RequireErrorHasCategory(&err, ErrorCategoryB(""))
				return errcat.Errorf(ErrZxcv, "aaah")
			}()
			return errcat.Recategorize(ErrQwer, err)
		}()
		shouldCategory(t, err, ErrQwer)
	})
	t.Run("assertions reject other categories of errors", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			fmt.Sprintf("...") // filler, to clarify line numbers
			if true {
				return errcat.Errorf(ErrZxcv, "aaah")
			}
			fmt.Sprintf("...") // filler, to clarify line numbers
			return nil
		}()
		shouldCategory(t, err, errcat.ErrCategoryFilterRejection)
		t.Logf("rejection for category'd errors:\n\t%s\n", err)
	})
	t.Run("assertions reject uncategorized errors", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			return fmt.Errorf("sad panda")
		}()
		shouldCategory(t, err, errcat.ErrCategoryFilterRejection)
		t.Logf("rejection for wild errors:\n\t%s\n", err)
	})
	t.Run("shadowing the func error is not a problem", func(t *testing.T) {
		err := func() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategoryA(""))
			if true {
				err := fmt.Errorf("sad panda")
				return err
			}
			return nil
		}()
		shouldCategory(t, err, errcat.ErrCategoryFilterRejection)
	})
}

func shouldCategory(t *testing.T, err error, cat interface{}) {
	t.Helper()
	ecat := errcat.Category(err)
	if ecat != cat {
		t.Errorf("expected category %v, got %v", cat, ecat)
	}
}
