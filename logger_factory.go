//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

package tlog

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	RegisterWriter(OutputConsole, DefaultConsoleWriterFactory)
	RegisterWriter(OutputFile, DefaultFileWriterFactory)
	Register(defaultLoggerName, NewZapLog(defaultConfig))
}

const (
	pluginType        = "log"
	defaultLoggerName = "default"
)

var (
	// DefaultLogger the default Logger. The initial output is console. When frame start, it is
	// over write by configuration.
	DefaultLogger Logger

	mu      sync.RWMutex
	loggers = make(map[string]Logger)
)

// Register registers Logger. It supports multiple Logger implementation.
func Register(name string, logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		panic("log: Register logger is nil")
	}
	if _, dup := loggers[name]; dup && name != defaultLoggerName {
		panic("log: Register called twiced for logger name " + name)
	}
	loggers[name] = logger
	if name == defaultLoggerName {
		DefaultLogger = logger
	}
}

// GetDefaultLogger gets the default Logger.
// To configure it, set key in configuration file to default.
// The console output is the default value.
func GetDefaultLogger() Logger {
	mu.RLock()
	l := DefaultLogger
	mu.RUnlock()
	return l
}

// SetLogger sets the default Logger.
func SetLogger(logger Logger) {
	mu.Lock()
	DefaultLogger = logger
	mu.Unlock()
}

// Get returns the Logger implementation by log name.
// log.Debug use DefaultLogger to print logs. You may also use log.Get("name").Debug.
func Get(name string) Logger {
	mu.RLock()
	l := loggers[name]
	mu.RUnlock()
	return l
}

// Sync syncs all registered loggers.
func Sync() {
	mu.RLock()
	defer mu.RUnlock()
	for _, logger := range loggers {
		_ = logger.Sync()
	}
}

// Decoder decodes the log.
type Decoder struct {
	OutputConfig *OutputConfig
	Core         zapcore.Core
	ZapLevel     zap.AtomicLevel
}

// Decode decodes writer configuration, copy one.
func (d *Decoder) Decode(cfg interface{}) error {
	output, ok := cfg.(**OutputConfig)
	if !ok {
		return fmt.Errorf("decoder config type:%T invalid, not **OutputConfig", cfg)
	}
	*output = d.OutputConfig
	return nil
}
