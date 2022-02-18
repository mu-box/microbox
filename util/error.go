package util

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/jcelliott/lumber"
)

type (
	Err struct {
		Code    string   // Code defining who is responsible for the error: 1xxx - user, 2xxx - hooks, 3xxx - images, 4xxx - platform, 5xxx - odin, 6xxx - cli
		Message string   // Error message
		Output  string   // Output from a command run
		Stack   []string // Origins of error
		Suggest string   // Suggested resolution
	}
)

// satisfy the error interface
func (eh Err) Error() string {
	if len(eh.Stack) == 0 {
		return eh.Message
	}
	return fmt.Sprintf("%s: %s", strings.Join(eh.Stack, ": "), eh.Message)
}

// log the error we ran into into our log file
func (eh Err) log() {
	// dont log if we are testing
	if flag.Lookup("test.v") != nil {
		return
	}

	lumber.Error(eh.Error())
	lumber.Error("%s\n", debug.Stack())
}

// Write an error message simular to Printf but logs the error to
// the log file
func ErrorfQuiet(fmtStr string, args ...interface{}) error {
	err := Err{
		Message: fmt.Sprintf(fmtStr, args...),
		Stack:   []string{},
	}
	err.log()
	return err
}

// Write an error message simular to Printf but logs the error to
// the log file
// todo: this is a silly workaround to preserve the suggestion
func ErrorfQuietErr(err error, args ...interface{}) error {
	newErr := Err{
		Message: fmt.Sprintf(err.Error(), args...),
		Stack:   []string{},
	}

	if err2, ok := err.(Err); ok {
		newErr.Suggest = err2.Suggest
		newErr.Output = err2.Output
		newErr.Code = err2.Code
	}

	newErr.log()
	return newErr
}

// creates an error the same fmt does
func Errorf(fmtStr string, args ...interface{}) error {
	err := Err{
		Message: fmt.Sprintf(fmtStr, args...),
		Stack:   []string{},
	}

	return err
}

// create an error
func ErrorQuiet(err error) error {
	if err == nil {
		return err
	}

	if er, ok := err.(Err); ok {
		return er
	}

	er := Err{
		Message: err.Error(),
		Stack:   []string{},
	}
	er.log()
	return er
}

// createson of our errors from a external error
func Error(err error) error {
	if err == nil {
		return err
	}

	eh := ErrorQuiet(err).(Err)
	return eh
}

// prepend the new message to the stack on our error messages
// this is useful because delimiting stack elements by :
// is not sufficient
func ErrorAppend(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(fmtStr, args...)

	// if it is one of our errors
	if er, ok := err.(Err); ok {
		// fmt.Println("OUR ERRORTYPE")
		er.Stack = append([]string{msg}, er.Stack...)
		return er
	}

	// make sure when we get any new error that isnt ours
	// we log it
	return ErrorAppend(Error(err), fmtStr, args...)
}
