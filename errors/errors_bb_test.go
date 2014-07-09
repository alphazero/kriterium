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

package errors_test

import (
	"github.com/elasticsearch/kriterium/errors"
	"testing"
	"testing/quick"
)

// ------------------------------------------------------------
// errors: black-box tests
// ------------------------------------------------------------

var quickConf = &quick.Config{MaxCount: 10000}

// quick check that TypedError#Code
func TestTypedError_Code(t *testing.T) {
	nopFn := func(code string) string {
		return code
	}
	testCodeFn := func(code string) string {
		e := errors.New(code)
		return e.Code()
	}
	fail := quick.CheckEqual(nopFn, testCodeFn, quickConf)
	if fail != nil {
		t.Error(fail)
	}
}

// quick check that TypedError#Matches
func TestTypedError_Intance(t *testing.T) {
	testCodeFn := func(code, extra string) bool {
		te := errors.New(code)
		e := te(extra)
		return te.Matches(e)
	}
	fail := quick.Check(testCodeFn, quickConf)
	if fail != nil {
		t.Error(fail)
	}
}
