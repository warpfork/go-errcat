package errcat

import (
	"fmt"
	"reflect"
	"runtime"
)

/*
	Filters an error value, forcing it to an ErrCategoryFilterRejection error if
	it does not have a category of the type specified.

	The typical/recommended usage for this is as a defer at the top of your
	function (so it's easy to see), bound to your actual return value (so it's
	impossible for an error to leave without hitting it):

		func foobar() (err error) {
			defer errcat.RequireErrorHasCategory(&err, ErrorCategory)
		}

	This makes for self-documenting code, and ensures that *if* you *do* make a
	coding error and return an inconsistent category, it is caught immediately --
	and we'll record the line number the error was returned from, so you can find
	and fix it quickly.

	(Yes, we all wish Go had a type system strong enough to simply check this at
	compile time, which is normal in other languages.  Alas.  Nonetheless, here's
	our attempt to do the best we can, even if it's merely at runtime.)

	This method mutates the error pointer you give it, so the error simple continues
	to return; it does not disrupt your control flow.
	You may also want to panic, though, since surely (surely; that's what you're
	declaring, if you use this feature) you are encountering a major bug: for this,
	use the `RequireErrorHasCategoryOrPanic` function.
*/
func RequireErrorHasCategory(e *error, category interface{}) {
	if err := requireErrorHasCategory(*e, category); err != nil {
		*e = err
	}
}

/*
	Identical to `RequireErrorHasCategory`, but panics.
*/
func RequireErrorHasCategoryOrPanic(e *error, category interface{}) {
	if err := requireErrorHasCategory(*e, category); err != nil {
		panic(err)
	}
}

func requireErrorHasCategory(e error, category interface{}) error {
	ecat := Category(e)
	switch ecat {
	case nil:
		return nil
	case ErrCategoryFilterRejection:
		// do nothing, because it's already redflagged.
		// (hm, or should we attach another line number?)
		return e
	case unknown:
		fallthrough
	default:
		rt_category := reflect.TypeOf(category)
		rt_ecat := reflect.TypeOf(ecat)
		if rt_ecat == rt_category {
			return nil
		}
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file, line = "?", 0
		}
		return ErrorDetailed(
			ErrCategoryFilterRejection,
			fmt.Sprintf("%s at %s:%d -- required %s, got %s(%q) (original error: %s)",
				ErrCategoryFilterRejection, file, line,
				rt_category.Name(), rt_ecat.Name(), ecat, e),
			Details(e),
		)
	}
}

const ErrCategoryFilterRejection = errorCategory("errcat-category-filter-rejection")
