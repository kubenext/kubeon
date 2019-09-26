/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrom(t *testing.T) {
	expected := NopLogger()
	ctx := WithLoggerContext(context.Background(), expected)
	actual := From(ctx)
	assert.True(t, actual == expected, "unexpected logger instance from context")
}

func TestMissingFrom(t *testing.T) {
	notExpected := NopLogger()
	ctx := context.Background()
	actual := From(ctx)
	assert.True(t, actual != notExpected, "unexpected logger instance from context")
	assert.NotNil(t, actual, "expected non-nil logger")
}

func BenchmarkFrom(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		From(context.Background())
	}
}
