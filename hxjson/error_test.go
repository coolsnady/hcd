// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hxjson_test

import (
	"testing"

	"github.com/coolsnady/hxd/hxjson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hxjson.ErrorCode
		want string
	}{
		{hxjson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{hxjson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{hxjson.ErrInvalidType, "ErrInvalidType"},
		{hxjson.ErrEmbeddedType, "ErrEmbeddedType"},
		{hxjson.ErrUnexportedField, "ErrUnexportedField"},
		{hxjson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{hxjson.ErrNonOptionalField, "ErrNonOptionalField"},
		{hxjson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{hxjson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{hxjson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{hxjson.ErrNumParams, "ErrNumParams"},
		{hxjson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(hxjson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hxjson.Error
		want string
	}{
		{
			hxjson.Error{Message: "some error"},
			"some error",
		},
		{
			hxjson.Error{Message: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
