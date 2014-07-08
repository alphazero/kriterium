// Licensed to Elasticsearch under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// package panics provides pseudo-exceptions.
package panics

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// -----------------------------------------------------------------------
// internal support
// -----------------------------------------------------------------------

// stringCodec allows for detecting types that provide a String() method
// internal use only.
type stringCodec interface {
	String() string
}

// -----------------------------------------------------------------------
// recoveredError with cause
// -----------------------------------------------------------------------

// recoveredError type with option cause.
// REVU: exporting this type and its member does not seem to be necessary.
type recoveredError struct {
	cause error
	err   error
}

func (e recoveredError) Error() string {
	return e.err.Error()
}

// Errors are returned by the panics package as plain 'error' references.
// This function is used to obtain of the underlying cause of such errors.
//
// If the argument is not a panics.recoveredError reference, then it simply returns
// the input argument.
func Cause(e error) error {
	ex, ok := e.(*recoveredError)
	if !ok {
		return e
	}
	return ex.cause
}

// -----------------------------------------------------------------------
// panics API
// -----------------------------------------------------------------------

// Asserts that input arg 'flag' is true.
// If false, panics with an recoveredError with descriptive cause based on the
// 'info' n-aray input arg.
func OnFalse(flag bool, info ...interface{}) {
	if flag {
		return
	}
	err := fmt.Errorf("%s - assert-fail:", fmtInfo(info...))
	panic(&recoveredError{cause: err, err: err}) // REVU: this dup use of 'err' is wrong
}

// Asserts that input arg 'flag' is false.
// If true, panics with an recoveredError with descriptive cause based on the
// 'info' n-aray input arg.
func OnTrue(flag bool, info ...interface{}) {
	if !flag {
		return
	}
	err := fmt.Errorf("%s - assert-fail:", fmtInfo(info...))
	panic(&recoveredError{cause: err, err: err}) // REVU: this dup use of 'err' is wrong
}

// Asserts that input arg 'v' is not nil.
// If nil, panics with an recoveredError with descriptive cause based on the
// 'info' n-aray input arg.
func OnNil(v interface{}, info ...interface{}) {
	if v != nil {
		return
	}
	err := fmt.Errorf("%s - value is nil:", fmtInfo(info...))
	panic(&recoveredError{cause: err, err: err}) // REVU: this dup use of 'err' is wrong
}

// Asserts that (error) input arg 'e' is nil.
// If not nil, panics with the input arg 'e'
// with descriptive cause based on the
// 'info' n-aray input arg.
func OnError(e error, info ...interface{}) {
	if e == nil {
		return
	}
	var err error = e
	if len(info) > 0 {
		err = fmt.Errorf("error: %s (cause: %s)", fmtInfo(info...), e)
	} else if !strings.HasPrefix(e.Error(), "error:") {
		err = fmt.Errorf("error: %s%s", fmtInfo(info...), e)
	}
	panic(&recoveredError{cause: e, err: err}) // REVU: this is correct use of error w/ cause
}

func fmtInfo(info ...interface{}) string {
	var msg = ""
	if len(info) > 0 {
		for _, s := range info {
			str := ""
			switch t := s.(type) {
			case string:
				str = t
			case stringCodec:
				str = t.String()
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				str = fmt.Sprintf("%d", t)
			case time.Time:
				str = fmt.Sprintf("'%d epoch-ns'", t.UnixNano())
			case bool:
				str = fmt.Sprintf("%t", t)
			default:
				str = fmt.Sprintf("%v", t)
			}
			str = " " + str
			msg += str
		}
		msg = strings.Trim(msg, " ")
	}
	return msg
}

// Recover encapsulates a generalized method of handing
// recovered panics, per std. panic/recover mechanism.
//
// Invocation of Recover() /must/ be deferred, per semantics of
// Go recover().
//
func Recover(err *error) error {
	if DEBUG {
		return nil
	}
	p := recover()
	if p == nil {
		return nil
	}

	switch t := p.(type) {
	case *recoveredError:
		*err = t
	case error:
		*err = t
	case string:
		*err = fmt.Errorf(t)
	default:
		*err = fmt.Errorf("recovered-panic: %q", t)
	}
	return *err
}

// AsyncRecover provides panic recovery utility analog to the
// panics.Recover function for go routines.
//
// Invocation of AsyncRecover() /must/ be deferred, per semantics of
// Go recover().
//
// Input args:
//
// stat: a user provided send-only channel for sending of
// recovery status. On recovery, either the recovered panic
// or the user provided okstat is sent on this channel. NOTE:
// channel len optimally should be of length 1.
//
// okstat: a user defined value used to signal that no panics
// occurred in the goroutine.
//
// TODO: no rush but refactor this ..
func AsyncRecover(stat chan<- interface{}, okstat interface{}) {
	if DEBUG {
		return
	}
	p := recover()
	if p == nil {
		stat <- okstat
		return
	}

	switch t := p.(type) {
	case *recoveredError:
		stat <- t
	case error:
		stat <- t
	case string:
		stat <- fmt.Errorf(t)
	default:
		stat <- fmt.Errorf("recovered-panic: %q", t)
	}
}

// Exist handler is analogous to Recover() but intended for top-level
// (e.g. main) functions that exit on error.
//
// Invocation of ExitHandler() /must/ be deferred, per semantics of
// Go recover().
//
// Input arg 'label' is purely informational and used in creation
// of the exit error.
func ExitHandler(label string) {
	if DEBUG {
		return
	}
	p := recover()
	if p == nil {
		os.Exit(0)
	}

	var e error
	switch t := p.(type) {
	case *recoveredError:
		e = t
	case error:
		e = t
	case string:
		e = fmt.Errorf(t)
	default:
		e = fmt.Errorf("recovered-panic: %q", t)
	}
	log.Fatalf("fatal error: %s: %s", label, e)
}

// set to true to short circuit the panic recovery mechanism
// and get the full stack dump per canonical panic().
var DEBUG = false

// ForFunc is a convenience feature that helps reduce code noise
// in functions that use panics API.
//
// For example, typically, we may want to include call-site
// informational bits that enhance the generated error per
// OnError, OnNil, etc. The repetition of providing these
// informational items adds no value. Defining these once
// at the start of the function/method addresses that.
//
// The returned panics.Panics interface mimics the top-level
// panics API.
//
// idiomatic use is to mask the package name so that there is
// minimal impact to existing code using e.g. panics.OnError, etc.
//
//    func something() (err error) {
//        defer panics.Recover(&err)
//        panics := panics.ForFunc("my-package/something():")
//        ...
//
//        stat, e := os.Stat("no-such-file")
//        panics.OnError(e)
//        ...
//
//    }
//
func ForFunc(fname string) Panics {
	return &fnpanics{fname}
}

type Panics interface {
	// See panics.Recover()
	Recover(err *error) error
	// See panics.OnError()
	OnError(e error, info ...interface{})
	// See panics.OnNil()
	OnNil(v interface{}, info ...interface{})
	// See panics.OnFalse()
	OnFalse(flag bool, info ...interface{})
	// See panics.OnTrue()
	OnTrue(flag bool, info ...interface{})
}

type fnpanics struct {
	fname string
}

func (t *fnpanics) Recover(err *error) error {
	e := Recover(err)
	return e
}

func (t *fnpanics) infoFixup(info ...interface{}) []interface{} {
	infofn := []interface{}{t.fname + "():"}
	return append(infofn, info...)
}
func (t *fnpanics) OnError(e error, info ...interface{}) {
	infofn := t.infoFixup(info...)
	log.Println(infofn)
	OnError(e, infofn...)
}
func (t *fnpanics) OnNil(v interface{}, info ...interface{}) {
	infofn := t.infoFixup(info...)
	OnNil(v, infofn...)
}
func (t *fnpanics) OnFalse(flag bool, info ...interface{}) {
	infofn := t.infoFixup(info...)
	OnFalse(flag, infofn...)
}
func (t *fnpanics) OnTrue(flag bool, info ...interface{}) {
	infofn := t.infoFixup(info...)
	OnTrue(flag, infofn...)
}
