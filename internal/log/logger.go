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

var contextKey = key("com.github.kubenext.kubeon")

type Logger interface {
	Debug(template string, args ...interface{})
	Info(template string, args ...interface{})
	Warn(template string, args ...interface{})
	Error(template string, args ...interface{})

	With(args ...interface{}) Logger
	WithErr(err error) Logger
	Named(name string) Logger
}

type sugaredLogWrapper struct {
	*zap.SugaredLogger
}

func (s *sugaredLogWrapper) Debug(template string, args ...interface{}) {
	s.SugaredLogger.Debugf(template, args)
}

func (s *sugaredLogWrapper) Info(template string, args ...interface{}) {
	s.SugaredLogger.Infof(template, args)
}

func (s *sugaredLogWrapper) Warn(template string, args ...interface{}) {
	s.SugaredLogger.Warnf(template, args)
}

func (s *sugaredLogWrapper) Error(template string, args ...interface{}) {
	s.Error(template, args)
}

func (s *sugaredLogWrapper) With(args ...interface{}) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With(args...)}
}

func (s *sugaredLogWrapper) WithErr(err error) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.With("err", err.Error())}
}

func (s *sugaredLogWrapper) Named(name string) Logger {
	return &sugaredLogWrapper{s.SugaredLogger.Named(name)}
}

var _ Logger = (*sugaredLogWrapper)(nil)

// Wrap zapSugaredLogger as Logger interface.
func Wrap(z *zap.SugaredLogger) Logger {
	return &sugaredLogWrapper{z}
}

// Return a new context a set logKey.
func WithLoggerContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
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
