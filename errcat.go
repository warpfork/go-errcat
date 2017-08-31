/*
	errcat is a simple universal error type that helps you produce
	errors that are both easy to categorize and handle, and also easy
	to maintain the original messages of.

	errcat does this by separating the two major parts of an error:
	the category and the message.

	The category is a value which you can switch on.
	*It is expected that the category field may be reassigned* as
	the error propagates up the stack.

	The message is the human-readable description of the error that occured.
	It *may* be further prepended with additional context info
	as it propagates out... or, not.
	The message may be redundant with the category: it is expected that
	the message will be printed to a user, while the category will
	not necessarily reach the user (it may be consumed by another layer
	of code, which may choose to re-categorize the error on its way up).

	Additional "details" may be attached in the Error.Details field;
	sometimes this can be used to provide key-value pairs which are
	useful in logging for other remote systems which must handle errors.
	However, usage of this should be minimized unless good reason is known;
	all handling logic should branch primarily on the category field,
	because that's what it's there for.

	errcat is specifically designed to be *serializable*, and just as
	importantly, *unserializable* again.
	This is helpful for making API-driven applications with
	consistent and reliably round-trip-able errors.
	errcat errors in json should appear as a very simple object:

		{"category":"your_tag", "msg":"full text goes here"}

	Typical usage patterns involve a const block in each package which
	enumerates the set of error category values that this package may return.
	When calling functions using the errcat convention, the callers may
	switch upon the returned Error's Category property:

		result, err := somepkg.SomeFunc()
		switch errcat.Category(err) {
		case nil:
			// good!  pass!
		case somepkg.ErrAlreadyDone:
			// good!  pass!
		case somepkg.ErrDataCorruption:
			// ... handle ...
		default:
			panic("bug: unknown error category")
		}

	Functions internal to packages may chose to panic up their errors.
	It is idiomatic to recover such internal panics and return the error
	as normal at the top of the package even when using panics as a
	non-local return system internally.
*/
package errcat
