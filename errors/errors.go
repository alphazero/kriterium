package errors

import "fmt"

// Error generator function type.
type ErrGenFn func(args ...interface{}) error

// Returns a new error generator function for the given error code.
//
// The generator function takes 0 or more generic arguments. Arguments
// are appended to the errcode parameter in the manner of fmt.Sprintln().
//
// If no args are provided, the generator function simply returns an error
// using the errcode provided and omits the ':' decoration after the errcode.
//
// Usage examples:
//
//    import "escargo/errors"
//    ...
//
//    var ErrWot   = errors.New("Something went wrong")
//    var ErrMyBad = errors.New("MyBad")
//    var ErrIO    = errors.New("IO Error")
//    ...
//
//    // ex: general error message with no additional details.
//    // always returns the error "ERR - Something went wrong"
//    func confusion() error {
//        if true {
//            return ErrWot()
//        }
//    }
//
//    // ex: general error code with instance specific parameters
//    // always returns the error "ERR - MyBad: oops, did it again :("
//    func errorProne() error {
//        if true {
//            return ErrMyBad("oops, did it again :(")
//        }
//    }
//
//    // ex: general error with underlying specific cause
//    // assuming that 'badarg' = "nosuchfile.txt"
//    // always returns the error "ERR - IOError: open nosuchfile.txt: no such file or directory"
//    func garbageIn(badarg string) error {
//        _, e := ioutils.ReadFile(badarg)
//        if e != nil {
//            return ErrIO(e.Error())
//        }
//    }
func New(errcode string) ErrGenFn {
	return func(args ...interface{}) error {
		decoration := ""
		if len(args) > 0 { decoration = ":"}
		errfmt := []interface{}{fmt.Sprintf("ERR - %s%s", errcode, decoration)}
		args0 := append(errfmt, args...)
		return fmt.Errorf(fmt.Sprintln(args0...))
	}
}
