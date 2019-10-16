// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xorcare/golden"
)

func TestEqual(t *testing.T) {
	testTree(t, func(t *testing.T) {
		golden.Equal(t, golden.Read(t)).FailNow()
		assert.NotNil(t, golden.Read(t), "input data cannot be empty")
		assert.NotNil(t, golden.SetTest(t).Read(), "golden data cannot be empty")
	})
}

func TestJSONEq(t *testing.T) {
	testTree(t, func(t *testing.T) {
		golden.JSONEq(t, string(golden.Read(t))).FailNow()
		assert.NotNil(t, golden.Read(t), "input data cannot be empty")
		assert.NotNil(t, golden.SetTest(t).SetPrefix("json").Read(), "golden data cannot be empty")
	})
}

// testTree needed to run a test function in tests with three levels of nesting.
func testTree(t *testing.T, f func(t *testing.T)) {
	// Simple test data without structure nesting.
	// testdata/TestExample.input
	// testdata/TestExample.golden
	f(t)
	t.Run("sublevel-one", func(t *testing.T) {
		// Tests with one level of nesting.
		// testdata/TestExample/sublevel-one.input
		// testdata/TestExample/sublevel-one.golden
		f(t)
		t.Run("sublevel-two", func(t *testing.T) {
			// Test with the second level of nesting.
			// testdata/TestExample/sublevel-one/sublevel-two.golden
			// testdata/TestExample/sublevel-one/sublevel-two.input
			f(t)
		})
	})
}
