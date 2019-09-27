/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestCreatePayload(t *testing.T) {
	action := "action"

	fields := map[string]interface{}{"foo": "bar"}
	got := CreatePayload(action, fields)

	expected := Payload{
		"action": action,
		"foo":    "bar",
	}
	assert.Equal(t, expected, got)
}

func TestPayload_GroupVersionKind(t *testing.T) {
	payload := Payload{
		"group":   "group",
		"version": "version",
		"kind":    "kind",
	}

	got, err := payload.GroupVersionKind()
	require.NoError(t, err)

	expected := schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "kind",
	}

	assert.Equal(t, expected, got)
}

func TestPayload_Float64(t *testing.T) {
	payload := Payload{
		"float": 64.66,
	}

	got, err := payload.Float64("float")
	require.NoError(t, err)

	expected := 64.66

	assert.Equal(t, expected, got)
}

func TestPayload_String(t *testing.T) {
	payload := Payload{
		"string": "string",
	}

	got, err := payload.String("string")
	require.NoError(t, err)

	expected := "string"

	assert.Equal(t, expected, got)
}

func TestPayload_OptionalString(t *testing.T) {
	payload := Payload{
		"string": "string",
	}

	got, err := payload.OptionalString("string")
	require.NoError(t, err)

	expected := "string"

	assert.Equal(t, expected, got)
}

func TestPayload_Uint16(t *testing.T) {
	tests := []struct {
		name     string
		payload  Payload
		key      string
		isErr    bool
		expected uint16
	}{
		{
			name:     "source is int",
			payload:  Payload{"uint16": float64(7)},
			key:      "uint16",
			expected: uint16(7),
		},
		{
			name:    "source overflows",
			payload: Payload{"uint16": 2 ^ 17},
			key:     "uint16",
			isErr:   true,
		},
		{
			name:    "source overflows",
			payload: Payload{"uint16": -1},
			key:     "uint16",
			isErr:   true,
		},
		{
			name:    "value is not int",
			payload: Payload{"uint16": true},
			key:     "uint16",
			isErr:   true,
		},
		{
			name:    "key does not exist",
			payload: Payload{},
			key:     "invalid",
			isErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.payload.Uint16(test.key)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, test.expected, got)
		})
	}
}
