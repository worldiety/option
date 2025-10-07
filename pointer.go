// Copyright (c) 2025 worldiety GmbH
//
// This file is part of the NAGO Low-Code Platform.
// Licensed under the terms specified in the LICENSE file.
//
// SPDX-License-Identifier: Custom-License

package option

import (
	"bytes"
	"encoding/json"
	"iter"
)

// Ptr is an immutable wrapper type around a pointer to a value of T. It returns always the T and never
// the pointer to T to express immutability. This is similar to Opt[T] but has an entire different
// memory characteristic. The optional implementation uses a flag and value type of T which requires
// the additional bool and the entire T. If you have a lot of options, e.g. in a nested struct, this will cause
// a lot of waste, but usually it lands on the stack.
//
// This Ptr implementation only keeps a pointer (usually 8 byte) to the heap escaped value but once given to
// the caller, the wrapped pointer cannot be mutated and only access to the dereferenced value is possible
// which creates a copy of T.
type Ptr[T any] struct {
	v *T
}

// Pointer is a factory to take the ownership of the given pointer.
func Pointer[T any](v *T) Ptr[T] {
	return Ptr[T]{
		v: v,
	}
}

// Unwrap makes the assertion that the Option is valid and otherwise panics. Such panic is always a programming error.
func (o Ptr[T]) Unwrap() T {
	return *o.v
}

// UnwrapOr returns either the contained value or the eagly evaluated given value. See also [Opt.UnwrapOrElse].
func (o Ptr[T]) UnwrapOr(v T) T {
	if o.IsNone() {
		return v
	}

	return *o.v
}

// UnwrapOrElse returns either the contained value or evaluates the given func to return an alternative.
func (o Ptr[T]) UnwrapOrElse(fn func() T) T {
	if o.IsNone() {
		return fn()
	}

	return *o.v
}

func (o Ptr[T]) IsZero() bool {
	return o.IsNone()
}

func (o Ptr[T]) IsSome() bool {
	return o.v != nil
}

func (o Ptr[T]) IsNone() bool {
	return o.v == nil
}

// All allows iteration over the possibly contained value.
func (o Ptr[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.v != nil {
			yield(*o.v)
		}
	}
}

func (o *Ptr[T]) UnmarshalJSON(buf []byte) error {
	if bytes.Equal([]byte("null"), buf) {
		o.v = nil
		return nil
	}

	var zero T // always allocate a new pointer otherwise we will break our immutability constraint for other copies
	err := json.Unmarshal(buf, &zero)
	if err != nil {
		return err
	}

	o.v = &zero
	return nil
}

func (o Ptr[T]) MarshalJSON() ([]byte, error) {
	if o.v != nil {
		return json.Marshal(o.v)
	}

	return []byte("null"), nil
}
