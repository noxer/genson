// Package genson provides helpers for encoding and decoding JSON values.
package genson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Any tries to decode a JSON field into the fields of C in order (T must be a struct type). It
// returns successfully after the first one that doesn't return an error. This is useful for
// interacting with APIs that return multiple different data types in the same field. After a
// successful decode, Field is set to the name of the field that as successfully decoded into.
type Any[T any] struct {
	FieldName string
	Payload   T
}

// MarshalJSON implements the JSON marshaler interface.
func (a *Any[T]) MarshalJSON() ([]byte, error) {
	val := reflect.ValueOf(&a.Payload).Elem()
	if val.Kind() != reflect.Struct {
		return nil, errors.New("must be a struct")
	}

	fields := reflect.VisibleFields(val.Type())

	if a.FieldName != "" {
		for _, field := range fields {
			if field.Name != a.FieldName {
				// skip non-matching fields
				continue
			}

			fval := val.FieldByIndex(field.Index)
			p, err := json.Marshal(fval.Interface())
			if err != nil {
				continue
			}

			return p, nil
		}
	} else {
		for _, field := range fields {
			fval := val.FieldByIndex(field.Index)
			if fval.IsZero() {
				// skip zero fields (0, nil, empty string)
				continue
			}

			p, err := json.Marshal(fval.Interface())
			if err != nil {
				continue
			}

			return p, nil
		}
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements the JSON unmarshaler interface.
func (a *Any[T]) UnmarshalJSON(data []byte) error {
	val := reflect.ValueOf(&a.Payload).Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("must be a struct")
	}

	var errs []error
	fields := reflect.VisibleFields(val.Type())
	for _, field := range fields {
		fval := val.FieldByIndex(field.Index)
		err := json.Unmarshal(data, fval.Addr().Interface())
		if err == nil {
			a.FieldName = field.Name
			return nil
		}

		errs = append(errs, fmt.Errorf("error encoding into %s: %w", field.Name, err))
	}

	return errors.Join(errs...)
}
