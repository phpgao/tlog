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
	"errors"
	"path/filepath"
)

var (
	// DefaultConsoleWriterFactory is the default console output implementation.
	DefaultConsoleWriterFactory = &ConsoleWriterFactory{}
	// DefaultFileWriterFactory is the default file output implementation.
	DefaultFileWriterFactory = &FileWriterFactory{}

	writers = make(map[string]PluginFactory)
)

type PluginFactory interface {
	// Type returns type of the plugin, i.e. selector, log, config, tracing.
	Type() string
	// Setup loads plugin by configuration.
	// The data structure of the configuration of the plugin needs to be defined in advance。
	Setup(name string, dec PluginDecoder) error
}

// PluginDecoder is the interface used to decode plugin configuration.
type PluginDecoder interface {
	Decode(cfg interface{}) error // the input param is the custom configuration of the plugin
}

// RegisterWriter registers log output writer. Writer may have multiple implementations.
func RegisterWriter(name string, writer PluginFactory) {
	writers[name] = writer
}

// GetWriter gets log output writer, returns nil if not exist.
func GetWriter(name string) PluginFactory {
	return writers[name]
}

// ConsoleWriterFactory is the console writer instance.
type ConsoleWriterFactory struct {
}

// Type returns the log plugin type.
func (f *ConsoleWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and registers console output writer.
func (f *ConsoleWriterFactory) Setup(name string, dec PluginDecoder) error {
	if dec == nil {
		return errors.New("console writer decoder empty")
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return errors.New("console writer log decoder type invalid")
	}
	cfg := &OutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	decoder.Core, decoder.ZapLevel = newConsoleCore(cfg)
	return nil
}

// FileWriterFactory is the file writer instance Factory.
type FileWriterFactory struct {
}

// Type returns log file type.
func (f *FileWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and register file output writer.
func (f *FileWriterFactory) Setup(name string, dec PluginDecoder) error {
	if dec == nil {
		return errors.New("file writer decoder empty")
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return errors.New("file writer log decoder type invalid")
	}
	if err := f.setupConfig(decoder); err != nil {
		return err
	}
	return nil
}

func (f *FileWriterFactory) setupConfig(decoder *Decoder) error {
	cfg := &OutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	if cfg.WriteConfig.LogPath != "" {
		cfg.WriteConfig.Filename = filepath.Join(cfg.WriteConfig.LogPath, cfg.WriteConfig.Filename)
	}
	if cfg.WriteConfig.RollType == "" {
		cfg.WriteConfig.RollType = RollBySize
	}

	core, level, err := newFileCore(cfg)
	if err != nil {
		return err
	}
	decoder.Core, decoder.ZapLevel = core, level
	return nil
}
