package option

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
)

// Opt is introduced because range over func can only represent at most 2 arguments. Processing
// a (T, ok, error) becomes impossible. Also, it is not correct to always use pointers for modelling or
// to use hidden error types for clear optional cases where an absent thing is never an error by definition.
// This also helps for performance edge cases, where you can technically express that a value is really
// just a value and does not escape.
//
// It sports also a non-nesting custom JSON serialization, which just uses NULL as representation.
// Note that if T is a pointer type, the Option becomes invalid after unmarshalling because a valid nil pointer
// cannot be distinguished from an invalid nil pointer in JSON, but you likely should not model your domain that
// way anyway.
//
// If you already have a pointer, just use its zero value which is nil and not Option.
// The zero value is safe to use.
type Opt[T any] struct {
	v     T
	valid bool
}

// Some is a factory to create a valid option.
func Some[T any](v T) Opt[T] {
	return Opt[T]{
		v:     v,
		valid: true,
	}
}

// None is only for better readability and equal to the zero value of Option.
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// Unwrap makes the assertion that the Option is valid and otherwise panics. Such panic is always a programming error.
func (o Opt[T]) Unwrap() T {
	if !o.valid {
		panic(fmt.Errorf("unwrapped invalid option"))
	}

	return o.v
}

// UnwrapOr returns either the contained value or the eagly evaluated given value. See also [Opt.UnwrapOrElse].
func (o Opt[T]) UnwrapOr(v T) T {
	if o.IsNone() {
		return v
	}

	return o.v
}

// UnwrapOrElse returns either the contained value or evaluates the given func to return an alternative.
func (o Opt[T]) UnwrapOrElse(fn func() T) T {
	if o.IsNone() {
		return fn()
	}

	return o.v
}

func (o Opt[T]) IsZero() bool {
	return o.IsNone()
}

func (o Opt[T]) IsSome() bool {
	return o.valid
}

func (o Opt[T]) IsNone() bool {
	return !o.valid
}

// All allows iteration over the possibly contained value.
func (o Opt[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.valid {
			yield(o.v)
		}
	}
}

func (o *Opt[T]) UnmarshalJSON(buf []byte) error {
	var zero T
	if bytes.Equal([]byte("null"), buf) {
		o.valid = false
		o.v = zero
		return nil
	}

	err := json.Unmarshal(buf, &zero)
	if err != nil {
		return err
	}

	o.valid = true
	o.v = zero
	return nil
}

func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if o.valid {
		return json.Marshal(o.v)
	}

	return []byte("null"), nil
}
