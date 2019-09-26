/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package strings

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestContains(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		sl       []string
		expected bool
	}{
		{
			name:     "does contain",
			s:        "1",
			sl:       []string{"1", "2", "3"},
			expected: true,
		},
		{
			name:     "does not contain",
			s:        "4",
			sl:       []string{"1", "2", "3"},
			expected: false,
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := Contains(tc.s, tc.sl)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDeduplicate(t *testing.T) {
	input := []string{"a", "a", "b", "c", "b"}
	got := Deduplicate(input)
	expected := []string{"a", "b", "c"}
	assert.Equal(t, expected, got)
}
