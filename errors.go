package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// OsExiter is the function used when the app exits. If not set defaults to os.Exit.
var OsExiter = os.Exit

// ErrWriter is used to write errors to the user. This can be anything
// implementing the io.Writer interface and defaults to os.Stderr.
var ErrWriter io.Writer = os.Stderr

// MultiError is an error that wraps multiple errors.
type MultiError interface {
	error
	Errors() []error
}

// newMultiError creates a new MultiError. Pass in one or more errors.
func newMultiError(err ...error) MultiError {
	ret := multiError(err)
	return &ret
}

type multiError []error

// Error implements the error interface.
func (m *multiError) Error() string {
	errs := make([]string, len(*m))
	for i, err := range *m {
		errs[i] = err.Error()
	}

	return strings.Join(errs, "\n")
}

// Errors returns a copy of the errors slice
func (m *multiError) Errors() []error {
	errs := make([]error, len(*m))
	copy(errs, *m)
	return errs
}

type requiredFlagsErr interface {
	error
}

type errRequiredFlags struct {
	missingFlags []string
}

func (e *errRequiredFlags) Error() string {
	if len(e.missingFlags) == 1 {
		return fmt.Sprintf("Required flag %q not set", e.missingFlags[0])
	}
	joinedMissingFlags := strings.Join(e.missingFlags, ", ")
	return fmt.Sprintf("Required flags %q not set", joinedMissingFlags)
}

type requiredArgumentsErr interface {
	error
}

type errRequiredArguments struct {
	missingArguments []string
}

func (e *errRequiredArguments) Error() string {
	if len(e.missingArguments) == 1 {
		return fmt.Sprintf("Required argument %q not set", e.missingArguments[0])
	}
	joinedMissingArguments := strings.Join(e.missingArguments, ", ")
	return fmt.Sprintf("Required arguments %q not set", joinedMissingArguments)
}

type mutuallyExclusiveGroup struct {
	flag1Name string
	flag2Name string
}

func (e *mutuallyExclusiveGroup) Error() string {
	return fmt.Sprintf("option %s cannot be set along with option %s", e.flag1Name, e.flag2Name)
}

type mutuallyExclusiveGroupRequiredFlag struct {
	flags *MutuallyExclusiveFlags
}

func (e *mutuallyExclusiveGroupRequiredFlag) Error() string {
	var missingFlags []string
	for _, grpf := range e.flags.Flags {
		var grpString []string
		for _, f := range grpf {
			grpString = append(grpString, f.Names()...)
		}
		missingFlags = append(missingFlags, strings.Join(grpString, " "))
	}

	return fmt.Sprintf("one of these flags needs to be provided: %s", strings.Join(missingFlags, ", "))
}

// ErrorFormatter is the interface that will suitably format the error output
type ErrorFormatter interface {
	Format(s fmt.State, verb rune)
}

// ExitCoder is the interface checked by `Command` for a custom exit code.
type ExitCoder interface {
	error
	ExitCode() int
}

type exitError struct {
	exitCode int
	err      error
}

// Exit wraps a message and exit code into an error, which by default is
// handled with a call to os.Exit during default error handling.
//
// This is the simplest way to trigger a non-zero exit code for a Command without
// having to call os.Exit manually. During testing, this behavior can be avoided
// by overriding the ExitErrHandler function on a Command or the package-global
// OsExiter function.
func Exit(message any, exitCode int) ExitCoder {
	var err error

	switch e := message.(type) {
	case ErrorFormatter:
		err = fmt.Errorf("%+v", message)
	case error:
		err = e
	default:
		err = fmt.Errorf("%+v", message)
	}

	return &exitError{
		err:      err,
		exitCode: exitCode,
	}
}

func (ee *exitError) Error() string {
	return ee.err.Error()
}

func (ee *exitError) ExitCode() int {
	return ee.exitCode
}

// Unwrap returns the error given to [Exit], so that [errors.Is] and
// [errors.As] can inspect it through the ExitCoder wrapper.
func (ee *exitError) Unwrap() error {
	return ee.err
}

// HandleExitCoder handles errors implementing ExitCoder by printing their
// message and calling OsExiter with the given exit code.
//
// If the given error instead implements MultiError, each error will be checked
// for the ExitCoder interface, and OsExiter will be called with the last exit
// code found, or exit code 1 if no ExitCoder is found.
//
// An ExitCoder that is wrapped with fmt.Errorf("%w") is handled the same way,
// so adding context to an error does not lose its exit code.
//
// This function is the default error-handling behavior for a Command.
func HandleExitCoder(err error) {
	if err == nil {
		return
	}

	exitErr, multiErr := findExitCoderOrMultiError(err)

	if exitErr != nil {
		if msg := err.Error(); msg != "" {
			if _, ok := exitErr.(ErrorFormatter); ok {
				_, _ = fmt.Fprintf(ErrWriter, "%+v\n", err)
			} else {
				_, _ = fmt.Fprintln(ErrWriter, err)
			}
		}
		OsExiter(exitErr.ExitCode())
		return
	}

	if multiErr != nil {
		code := handleMultiError(multiErr)
		OsExiter(code)
		return
	}
}

// findExitCoderOrMultiError walks the wrap chain of err and returns the first
// error that is an ExitCoder or a MultiError, whichever comes first. Unlike
// [errors.As], it stops at a MultiError instead of descending into it (a
// MultiError may also implement Unwrap() []error), so that the "last exit code
// wins" rule of MultiError is kept.
func findExitCoderOrMultiError(err error) (ExitCoder, MultiError) {
	for err != nil {
		if exitErr, ok := err.(ExitCoder); ok {
			return exitErr, nil
		}
		if multiErr, ok := err.(MultiError); ok {
			return nil, multiErr
		}

		switch u := err.(type) {
		case interface{ Unwrap() error }:
			err = u.Unwrap()
		case interface{ Unwrap() []error }:
			for _, e := range u.Unwrap() {
				if exitErr, multiErr := findExitCoderOrMultiError(e); exitErr != nil || multiErr != nil {
					return exitErr, multiErr
				}
			}
			return nil, nil
		default:
			return nil, nil
		}
	}
	return nil, nil
}

func handleMultiError(multiErr MultiError) int {
	code := 1
	for _, merr := range multiErr.Errors() {
		if multiErr2, ok := merr.(MultiError); ok {
			code = handleMultiError(multiErr2)
		} else if merr != nil {
			fmt.Fprintln(ErrWriter, merr)
			if c, ok := exitCodeOf(merr); ok {
				code = c
			}
		}
	}
	return code
}

// exitCodeOf reports the exit code carried by err, which may be wrapped: the
// code of an ExitCoder, or the last exit code found in a MultiError.
func exitCodeOf(err error) (int, bool) {
	exitErr, multiErr := findExitCoderOrMultiError(err)
	if exitErr != nil {
		return exitErr.ExitCode(), true
	}
	if multiErr != nil {
		code, found := 0, false
		for _, e := range multiErr.Errors() {
			if c, ok := exitCodeOf(e); ok {
				code, found = c, true
			}
		}
		return code, found
	}
	return 0, false
}
