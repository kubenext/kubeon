/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package log

import (
	"context"
	"go.uber.org/zap"
)

type key string

var contextKey = key("com.kubenext.logger")

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})

	With(args ...interface{}) Logger
	WithErr(err error) Logger
	Named(name string) Logger
}

type sugaredLoggerWrapper struct {
	*zap.SugaredLogger
}

func (s *sugaredLoggerWrapper) With(args ...interface{}) Logger {
	panic("implement me")
}

func (s *sugaredLoggerWrapper) WithErr(err error) Logger {
	panic("implement me")
}

func (s *sugaredLoggerWrapper) Named(name string) Logger {
	panic("implement me")
}

var _ Logger = (*sugaredLoggerWrapper)(nil)

func Wrap(zap *zap.SugaredLogger) Logger {
	return &sugaredLoggerWrapper{zap}
}

func NopLogger() Logger {
	return Wrap(zap.NewNop().Sugar())
}

func From(ctx context.Context) Logger {
	if ctx == nil {
		return NopLogger()
	}

	v := ctx.Value(contextKey)
	log, ok := v.(Logger)
	if !ok || log == nil {
		return NopLogger()
	}
	return log
}

func WithLoggerContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}
