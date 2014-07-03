// basic usage example of pacakge escargo/errors
package main

import (
	"log"
	"flag"
	"io/ioutil"
	"escargo/errors"
)

// This form is only a suggestion. We find it useful and clean.
var ERR = struct {
		IllegalArgument, IOError, FileNotFound, PemDecodeError, MyBadError errors.ErrGenFn
	}{
	IllegalArgument: errors.New("IllegalArgument"),
	IOError:         errors.New("IOError"),
	FileNotFound:    errors.New("FileNotFound"),
}

// alternatively you can just simply define our vars as a var.
var ErrorXYZ = errors.New("Error XYZ")
	
var options = struct {
		filepath string
	}{
	filepath: "foo.bar",
}

func init() {
	log.SetFlags(0)
	flag.StringVar(&options.filepath, "f", options.filepath, "path to file to read")
}

func main() {
	flag.Parse()
	log.Println("filepath:", options.filepath)

	e := loadFile(options.filepath)
	if e != nil {
		log.Fatalf(e.Error())
	}
}

func loadFile(path string) error {
	if path == "" {
		return ERR.IllegalArgument("path is nil")
	}
	_, e := ioutil.ReadFile(path)
	if e != nil {
		return ERR.IOError(e.Error())
	}
	return ErrorXYZ ("oops, did it again :(")
}
