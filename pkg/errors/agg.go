package errors

import (
	"errors"
	"gim/pkg/store"
)

// Aggregate represents an object that contains multiple errors, but does not
// necessarily have singular semantic meaning.
// The aggregate can be used with `errors.Is()` to check for the occurrence of
// a specific error type.
// Errors.As() is not supported, because the caller presumably cares about a
// specific error of potentially multiple that match the given type.
type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

// NewAggregate converts a slice of errors into an Aggregate interface, which
// is itself an implementation of the error interface.  If the slice is empty,
// this returns nil.
// It will check if any of the element of input error list is nil, to avoid
// nil pointer panic when call Error().
func NewAggregate(errlist []error) Aggregate {
	if len(errlist) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return aggregate(errs)
}

// This helper implements the error and Errors interfaces.  Keeping it private
// prevents people from making an aggregate of 0 errors, which is not
// an error, but does satisfy the error interface.
type aggregate []error

// Error is part of the error interface.
func (agg aggregate) Error() string {
	if len(agg) == 0 {
		// This should never happen, really.
		return ""
	}
	if len(agg) == 1 {
		return agg[0].Error()
	}
	seenerrs := store.ThreadSafe()
	result := ""
	agg.visit(func(err error) bool {
		msg := err.Error()
		if seenerrs.Contain(msg) {
			return false
		}
		seenerrs.Add(msg)
		if seenerrs.Size()> 1 {
			result += ", "
		}
		result += msg
		return false
	})
	if seenerrs.Size() > 1 {
		return result
	}
	return "[" + result + "]"
}

func (agg aggregate) Is(target error) bool {
	return agg.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

func (agg aggregate) visit(f func(err error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate:
			if match := err.visit(f); match {
				return match
			}
		case Aggregate:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

// Errors is part of the Aggregate interface.
func (agg aggregate) Errors() []error {
	return []error(agg)
}
