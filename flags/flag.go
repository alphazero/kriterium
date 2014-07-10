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

// Enhanced command line flag definition and utility.
//
// Usage example (main command line options):
//		package main
//
//		import (
//			"flag"
//			"github.com/elasticsearch/kriterium/flags"
//			"github.com/elasticsearch/kriterium/panics"
//			"log"
//		)
//
//		var options = &struct {
//				Global   *flags.BoolOption
//				Id       *flags.StringOption
//				LongOnly *flags.IntOption
//			}{
//			Global:   flags.NewBoolOption(flag.CommandLine, "g", "global", false, "optional boolean flag", false), // note final false arg
//			Id:       flags.NewStringOption(flag.CommandLine, "id", "identity", "", "required flag with both long and short name", true),
//			LongOnly: flags.NewIntOption(flag.CommandLine, "", "just-long", 100, "required int flag with a long name only", false),
//		}
//
//		func main() {
//			defer panics.ExitHandler("command-line options example")
//
//			// parse
//			flag.Parse()
//
//			// check
//			e := flags.UsageVerify(options)
//			panics.OnError(e)
//
//			// access
//			log.SetFlags(0)
//			log.Printf("option: %s default:%v provided:%v\n", options.Global.LongName(), options.Global.Default(), options.Global.Get())
//			log.Printf("option: %s default:%v provided:%v\n", options.Id.LongName(), options.Id.Default(), options.Id.Get())
//			log.Printf("option: %s default:%v provided:%v\n", options.LongOnly.LongName(), options.LongOnly.Default(), options.LongOnly.Get())
//		}
//
// [todo: usage example for sub-commands]
//
package flags

import (
	"flag"
	"fmt"
	"github.com/elasticsearch/kriterium/errors"
	"reflect"
)

// -----------------------------------------------------------------------
// command line options and flags
// -----------------------------------------------------------------------
type optionSpec struct {
	short, long, info string
	required          bool
}

type Option interface {
	Kind() reflect.Kind
	Name() string
	LongName() string
	Provided() bool
	Required() bool
	Info() string
}

func verifyRequiredOption(option Option) error {
	if option == nil {
		return errors.IllegalArgument("verifyRequiredOption:", "option is nil")
	}
	if option.Provided() {
		return nil
	}
	shortIfAny := option.Name()
	longIfAny := option.LongName()
	var usage string
	switch {
	case shortIfAny == "" && longIfAny == "":
		panic(errors.IllegalState("BUG:", "option:", option, "neither short or long name defined"))
	case shortIfAny == "":
		usage = fmt.Sprintf("Option { -%s } must be provided", longIfAny)
	case longIfAny == "":
		usage = fmt.Sprintf("Option { -%s } must be provided", shortIfAny)
	default:
		usage = fmt.Sprintf("Option { -%s | -%s } must be provided", shortIfAny, longIfAny)
	}
	return errors.RequiredFlag(usage)
}

func UsageVerify(str interface{}) error {
	s := reflect.ValueOf(str).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		option := f.Interface().(Option)
		if option.Required() {
			if e := verifyRequiredOption(option); e != nil {
				return e
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////
/// generated code - do not edit ///////////////////////////////////
////////////////////////////////////////////////////////////////////

/// bool //////////////////////////////////////////////////

type BoolOption struct {
	optionSpec
	value  *bool
	defval bool
}

type BoolValue interface {
	Get() bool
	Default() bool
}

func (opt *BoolOption) Get() bool {
	return *opt.value
}

func (opt *BoolOption) Default() bool {
	return opt.defval
}

func (t *BoolOption) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.BoolVar(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.BoolVar(t.value, t.long, t.defval, t.info)
	}
}

func (opt *BoolOption) Kind() reflect.Kind {
	return reflect.Bool
}

func (opt *BoolOption) Name() string {
	return opt.short
}

func (opt *BoolOption) LongName() string {
	return opt.long
}

func (opt *BoolOption) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *BoolOption) Required() bool {
	return opt.required
}

func (opt *BoolOption) Info() string {
	return opt.info
}

func newBoolOption(short, long, info string, defval bool, required bool) *BoolOption {
	return &BoolOption{optionSpec{short, long, info, required}, new(bool), defval}
}

// TODO: rename to ..Option
func NewBoolOption(fs *flag.FlagSet, short, long string, defval bool, info string, required bool) *BoolOption {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newBoolOption(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// int //////////////////////////////////////////////////

type IntOption struct {
	optionSpec
	value  *int
	defval int
}

type IntValue interface {
	Get() int
	Default() int
}

func (opt *IntOption) Get() int {
	return *opt.value
}

func (opt *IntOption) Default() int {
	return opt.defval
}

func (t *IntOption) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.IntVar(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.IntVar(t.value, t.long, t.defval, t.info)
	}
}

func (opt *IntOption) Kind() reflect.Kind {
	return reflect.Int
}

func (opt *IntOption) Name() string {
	return opt.short
}

func (opt *IntOption) LongName() string {
	return opt.long
}

func (opt *IntOption) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *IntOption) Required() bool {
	return opt.required
}

func (opt *IntOption) Info() string {
	return opt.info
}

func newIntOption(short, long, info string, defval int, required bool) *IntOption {
	return &IntOption{optionSpec{short, long, info, required}, new(int), defval}
}

// TODO: rename to ..Option
func NewIntOption(fs *flag.FlagSet, short, long string, defval int, info string, required bool) *IntOption {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newIntOption(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// int64 //////////////////////////////////////////////////

type Int64Option struct {
	optionSpec
	value  *int64
	defval int64
}

type Int64Value interface {
	Get() int64
	Default() int64
}

func (opt *Int64Option) Get() int64 {
	return *opt.value
}

func (opt *Int64Option) Default() int64 {
	return opt.defval
}

func (t *Int64Option) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.Int64Var(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.Int64Var(t.value, t.long, t.defval, t.info)
	}
}

func (opt *Int64Option) Kind() reflect.Kind {
	return reflect.Int64
}

func (opt *Int64Option) Name() string {
	return opt.short
}

func (opt *Int64Option) LongName() string {
	return opt.long
}

func (opt *Int64Option) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *Int64Option) Required() bool {
	return opt.required
}

func (opt *Int64Option) Info() string {
	return opt.info
}

func newInt64Option(short, long, info string, defval int64, required bool) *Int64Option {
	return &Int64Option{optionSpec{short, long, info, required}, new(int64), defval}
}

// TODO: rename to ..Option
func NewInt64Option(fs *flag.FlagSet, short, long string, defval int64, info string, required bool) *Int64Option {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newInt64Option(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// uint //////////////////////////////////////////////////

type UintOption struct {
	optionSpec
	value  *uint
	defval uint
}

type UintValue interface {
	Get() uint
	Default() uint
}

func (opt *UintOption) Get() uint {
	return *opt.value
}

func (opt *UintOption) Default() uint {
	return opt.defval
}

func (t *UintOption) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.UintVar(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.UintVar(t.value, t.long, t.defval, t.info)
	}
}

func (opt *UintOption) Kind() reflect.Kind {
	return reflect.Uint
}

func (opt *UintOption) Name() string {
	return opt.short
}

func (opt *UintOption) LongName() string {
	return opt.long
}

func (opt *UintOption) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *UintOption) Required() bool {
	return opt.required
}

func (opt *UintOption) Info() string {
	return opt.info
}

func newUintOption(short, long, info string, defval uint, required bool) *UintOption {
	return &UintOption{optionSpec{short, long, info, required}, new(uint), defval}
}

// TODO: rename to ..Option
func NewUintOption(fs *flag.FlagSet, short, long string, defval uint, info string, required bool) *UintOption {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newUintOption(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// uint64 //////////////////////////////////////////////////

type Uint64Option struct {
	optionSpec
	value  *uint64
	defval uint64
}

type Uint64Value interface {
	Get() uint64
	Default() uint64
}

func (opt *Uint64Option) Get() uint64 {
	return *opt.value
}

func (opt *Uint64Option) Default() uint64 {
	return opt.defval
}

func (t *Uint64Option) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.Uint64Var(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.Uint64Var(t.value, t.long, t.defval, t.info)
	}
}

func (opt *Uint64Option) Kind() reflect.Kind {
	return reflect.Uint64
}

func (opt *Uint64Option) Name() string {
	return opt.short
}

func (opt *Uint64Option) LongName() string {
	return opt.long
}

func (opt *Uint64Option) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *Uint64Option) Required() bool {
	return opt.required
}

func (opt *Uint64Option) Info() string {
	return opt.info
}

func newUint64Option(short, long, info string, defval uint64, required bool) *Uint64Option {
	return &Uint64Option{optionSpec{short, long, info, required}, new(uint64), defval}
}

// TODO: rename to ..Option
func NewUint64Option(fs *flag.FlagSet, short, long string, defval uint64, info string, required bool) *Uint64Option {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newUint64Option(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// float64 //////////////////////////////////////////////////

type Float64Option struct {
	optionSpec
	value  *float64
	defval float64
}

type Float64Value interface {
	Get() float64
	Default() float64
}

func (opt *Float64Option) Get() float64 {
	return *opt.value
}

func (opt *Float64Option) Default() float64 {
	return opt.defval
}

func (t *Float64Option) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.Float64Var(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.Float64Var(t.value, t.long, t.defval, t.info)
	}
}

func (opt *Float64Option) Kind() reflect.Kind {
	return reflect.Float64
}

func (opt *Float64Option) Name() string {
	return opt.short
}

func (opt *Float64Option) LongName() string {
	return opt.long
}

func (opt *Float64Option) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *Float64Option) Required() bool {
	return opt.required
}

func (opt *Float64Option) Info() string {
	return opt.info
}

func newFloat64Option(short, long, info string, defval float64, required bool) *Float64Option {
	return &Float64Option{optionSpec{short, long, info, required}, new(float64), defval}
}

// TODO: rename to ..Option
func NewFloat64Option(fs *flag.FlagSet, short, long string, defval float64, info string, required bool) *Float64Option {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newFloat64Option(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

/// string //////////////////////////////////////////////////

type StringOption struct {
	optionSpec
	value  *string
	defval string
}

type StringValue interface {
	Get() string
	Default() string
}

func (opt *StringOption) Get() string {
	return *opt.value
}

func (opt *StringOption) Default() string {
	return opt.defval
}

func (t *StringOption) defineFlag(fs *flag.FlagSet) {
	if t.short != "" {
		fs.StringVar(t.value, t.short, t.defval, t.info)
	}
	if t.long != "" {
		fs.StringVar(t.value, t.long, t.defval, t.info)
	}
}

func (opt *StringOption) Kind() reflect.Kind {
	return reflect.String
}

func (opt *StringOption) Name() string {
	return opt.short
}

func (opt *StringOption) LongName() string {
	return opt.long
}

func (opt *StringOption) Provided() bool {
	return *opt.value != opt.defval
}

func (opt *StringOption) Required() bool {
	return opt.required
}

func (opt *StringOption) Info() string {
	return opt.info
}

func newStringOption(short, long, info string, defval string, required bool) *StringOption {
	return &StringOption{optionSpec{short, long, info, required}, new(string), defval}
}

// TODO: rename to ..Option
func NewStringOption(fs *flag.FlagSet, short, long string, defval string, info string, required bool) *StringOption {
	if short == "" && long == "" {
		panic(errors.IllegalArgument("Either long or short flag name must be provided"))
	} else if short == long {
		panic(errors.IllegalArgument("flag long & short names must be distinct"))
	}
	opt := newStringOption(short, long, info, defval, required)
	opt.defineFlag(fs)
	return opt
}

////////////////////////////////////////////////////////////////////
/// generated code - do not edit /////////////////////////// end ///
////////////////////////////////////////////////////////////////////
