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

package panics_test

import (
	"kriterium/panics"
	"testing"
	//	"testing/quick"
	"fmt"
)

// TODO: quickcheck all params
//var quickConf = &quick.Config{MaxCount: 10000}

// used to check callsite handling of error returning func
var okFunc = func() error { return nil }

// used to check callsite handling of error returning func
// all functions defined to use panics.On???()
var errFuncsApi = []func() error{
	func() (err error) {
		defer panics.Recover(&err)
		e := fmt.Errorf("test-error")
		panics.OnError(e)
		return
	},
	func() (err error) {
		defer panics.Recover(&err)
		panics.OnNil(nil)
		return
	},
	func() (err error) {
		defer panics.Recover(&err)
		panics.OnFalse(false)
		return
	},
	func() (err error) {
		defer panics.Recover(&err)
		panics.OnTrue(true)
		return
	},
}

// used to check callsite handling of error returning func
// all functions defined to use panics.ForFunc() & panics.On???()
var errFuncsForFunc = []func() error{
	func() (err error) {
		panics := panics.ForFunc("errFuncsForFunc[onError]")
		defer panics.Recover(&err)
		e := fmt.Errorf("test-error")
		panics.OnError(e)
		return
	},
	func() (err error) {
		panics := panics.ForFunc("errFuncsForFunc[OnNil]")
		defer panics.Recover(&err)
		panics.OnNil(nil)
		return
	},
	func() (err error) {
		panics := panics.ForFunc("errFuncsForFunc[onFalse]")
		defer panics.Recover(&err)
		panics.OnFalse(false)
		return
	},
	func() (err error) {
		panics := panics.ForFunc("errFuncsForFunc[onTrue]")
		defer panics.Recover(&err)
		panics.OnTrue(true)
		return
	},
}

// test panics API - sync
func TestPanicsErrorsAndRecoverAPI(t *testing.T) {
	for _, fn := range errFuncsApi {
		e := fn()
		if e == nil {
			t.Error("expected to return error")
		}
	}
}

// test panics.ForFunc - sync
func TestPanicsErrorsAndRecoverForFunc(t *testing.T) {
	for _, fn := range errFuncsForFunc {
		e := fn()
		if e == nil {
			t.Error("expected to return error")
		}
	}
}

// test panics API functionality at call site - async
// functions tested use panics API
func TestAsyncPanicsErrorsAndRecoverForAPI(t *testing.T) {
	var okStat = "ok"

	asyncFn := func(fn func() error, statchan chan<- interface{}) {
		defer panics.AsyncRecover(statchan, okStat)
		panics.OnError(fn())
	}
	statchan := make(chan interface{}, 1)

	for _, fn := range errFuncsApi {
		go asyncFn(fn, statchan)

		stat := <-statchan
		if stat == nil {
			t.Error("AsyncRecover should never return nil")
		}
		if stat == okStat {
			t.Error("expected error not okStat")
		}
		_, ok := stat.(error)
		if !ok {
			t.Error("expected error type statt")
		}
	}

	// now test that okStat is correctly returned in case of no errors
	go asyncFn(okFunc, statchan)

	stat := <-statchan
	if stat == nil {
		t.Error("AsyncRecover should never return nil")
	}
	if stat != okStat {
		t.Error("expected okStat")
	}
}

// test panics API functionality at call site - async
// functions tested use panics.ForFunc API
func TestAsyncPanicsErrorsAndRecoverForFunc(t *testing.T) {
	var okStat = "ok"

	asyncFn := func(fn func() error, statchan chan<- interface{}) {
		defer panics.AsyncRecover(statchan, okStat)
		panics.OnError(fn())
	}
	statchan := make(chan interface{}, 1)

	for _, fn := range errFuncsForFunc {
		go asyncFn(fn, statchan)

		stat := <-statchan
		if stat == nil {
			t.Error("AsyncRecover should never return nil")
		}
		if stat == okStat {
			t.Error("expected error not okStat")
		}
		_, ok := stat.(error)
		if !ok {
			t.Error("expected error type statt")
		}
	}

	// now test that okStat is correctly returned in case of no errors
	go asyncFn(okFunc, statchan)

	stat := <-statchan
	if stat == nil {
		t.Error("AsyncRecover should never return nil")
	}
	if stat != okStat {
		t.Error("expected okStat")
	}
}
