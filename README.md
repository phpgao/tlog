# tlog: The log module extracted from trpc-group/trpc-go

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/phpgao/tlog/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/phpgao/tlog?style=social)](https://github.com/phpgao/tlog)

tlog is a lightweight and flexible logging library for Go applications, extracted from the trpc-group/trpc-go project. It provides a simple way to configure and manage logs in your applications with various output configurations and levels.

## Features

- **Multiple Output Configurations**: Support for different writers like console and file.
- **Log Level Management**: Easily set log levels for different parts of your application.
- **Custom Formatters**: JSON formatting support for structured logging.
- **Contextual Logging**: Add context to your logs for better traceability.

## Installation

To get started with tlog, simply install it using go get:

```bash
go get github.com/phpgao/tlog
```

## Configuration

tlog allows you to configure multiple output destinations and set log levels and formatters for each. Here's an example configuration:

```go
package main

import (
    "github.com/phpgao/tlog"
    "net/http"
)

func main() {
    var c = []tlog.OutputConfig{
        {
            Writer:    "console",
            Level:     "info",
            Formatter: "json",
        },
        {
            Writer: "file",
            WriteConfig: tlog.WriteConfig{
                LogPath:    "/tmp/",
                Filename:   "trpc.log",
                RollType:   "size",
                MaxAge:     1,
                MaxBackups: 5,
                Compress:   false,
                MaxSize:    1024,
            },
            Level:     "debug",
            Formatter: "json",
        },
    }

    tlog.Register("default", tlog.NewZapLog(c))
    tlog.Infof("Hello, tlog!")

	tlog.SetLevel("0", tlog.LevelInfo)
}
```

## Usage

tlog provides various methods to log messages at different levels. Here's how you can use it:

```go
// Log an informational message
tlog.Infof("Info message: %s", "This is an info message")

// Log an error message
tlog.Errorf("Error message: %s", "This is an error message")

// Log a message with context
ctx := tlog.WithContextFields(context.TODO(), "key1", "value1", "key2", "value2")
tlog.InfoContextf(ctx, "Contextual message: %s", "This message has context")

// use RegisterHandler to set log level
mux := http.NewServeMux()
RegisterHandlerWithPath(mux, "/")
http.ListenAndServe(":8080", mux)

// gin logger
router.Use(handler.GinLogger(), gin.Recovery())
```

## License

tlog is released under the [Apache License Version 2.0](https://github.com/phpgao/tlog/blob/main/LICENSE.txt).

---

Feel free to star and watch the repository to keep up with the latest updates!
```