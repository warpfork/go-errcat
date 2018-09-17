// Typed nils are *extremely* annoying.
//
// Typed nils have a very high likelihood in practical code as people first
// draft it of causing nil dereference errors to rise *from in our library*
// when someone implements e.g. the 'Category()' function without checking
// for their own nilness.  This is ugly and useless distraction.
//
// So.  How expensive is it to add a reflect-based check that "does the right
// thing" for typed nils?
//
package nilly

import (
	"reflect"
	"testing"

	"github.com/warpfork/go-errcat"
)

// Ballpark results:
//
//		BenchmarkNilCheck-8                     2000000000               0.27 ns/op            0 B/op          0 allocs/op
//		BenchmarkReflectNilCheck-8              300000000                3.98 ns/op            0 B/op          0 allocs/op
//		BenchmarkReflectNilCheckCorrectly-8     500000000                3.50 ns/op            0 B/op          0 allocs/op
//
// Twelve times slower.  Unfortunately significant.
//
// And further, this actually gets worse if you use a struct rather than pointer,
// because now you'll get allocs in there somewhere:
//
//		BenchmarkNilCheck_Nonptr-8                      2000000000               0.27 ns/op            0 B/op          0 allocs/op
//		BenchmarkReflectNilCheckCorrectly_Nonptr-8      30000000                49.1 ns/op            48 B/op          1 allocs/op
//
// Conclusion: this is basically a no-go if we want this library to be near-zero cost.
//

// Related learning: changing all of the interface methods to only work on the
// struct rather than the pointer does at least make somewhat more sensible
// error messages:
//
//  "panic: value method nilly.errStruct.Category called using nil *errStruct pointer"
//
// However, overall, this is not particularly helpful to our ergonomics in this
// library, because though it gives a clearer message, people will still tend
// to get into the situation where it's a problem: it's typical to write
// structs with a typed error pointer as a way of giving a concrete type hint
// to serialize/deserialize systems and simultaneously having a clear way to
// indicate 'this is zero, plz do not serialize an empty object here'.
//

func BenchmarkNilCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ex := exampleContainer{}
		err := errcat.Error(ex.Error)
		if err == nil {
			panic("typed nils don't work this way")
		}
	}
}

func BenchmarkReflectNilCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ex := exampleContainer{}
		err := errcat.Error(ex.Error)
		if reflect.ValueOf(err).IsNil() == false {
			panic("reflect should have caught this")
		}
	}
}

func BenchmarkReflectNilCheckCorrectly(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ex := exampleContainer{}
		err := errcat.Error(ex.Error)
		rv := reflect.ValueOf(err)
		if rv.Kind() == reflect.Ptr && rv.IsNil() == false {
			panic("reflect should have caught this")
		}
	}
}

type exampleContainer struct {
	Error *errStruct
}

type errStruct struct {
	Category_ interface{}       `json:"category"          refmt:"category"`
	Message_  string            `json:"message"           refmt:"message"`
	Details_  map[string]string `json:"details,omitempty" refmt:"details,omitempty"`
}

func (e errStruct) Category() interface{}      { return e.Category_ }
func (e errStruct) Message() string            { return e.Message_ }
func (e errStruct) Details() map[string]string { return e.Details_ }
func (e errStruct) Error() string              { return e.Message_ }

func BenchmarkNilCheck_Nonptr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ex := exampleContainer2{}
		err := errcat.Error(ex.Error)
		if err == nil {
			panic("typed nils don't work this way")
		}
	}
}

func BenchmarkReflectNilCheck_Nonptr(b *testing.B) {
	b.Skip("this would panic because reflect.Value.IsNil() is rude")
}

func BenchmarkReflectNilCheckCorrectly_Nonptr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ex := exampleContainer2{}
		err := errcat.Error(ex.Error)
		rv := reflect.ValueOf(err)
		if rv.Kind() == reflect.Ptr && rv.IsNil() == false {
			panic("reflect should have caught this")
		}
	}
}

// all the same except *not* a pointer to errStruct
type exampleContainer2 struct {
	Error errStruct
}

//func BenchmarkSanityCheck(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		ex := exampleContainer{}
//		if errcat.Category(ex.Error) == nil {
//			panic("typed nils don't work this way")
//		}
//	}
//}
