// Package handler comes from admin/admin.go of trpc-go
// you can add this handler to your http(s) server
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/phpgao/tlog"
)

// PattenLoglevel View/Set the log level of the framework
const PattenLoglevel = "/cmds/loglevel"

// return param.
var (
	ReturnErrCodeParam = "errorcode"
	ReturnMessageParam = "message"
	ErrCodeServer      = 1
)

func RegisterHandler(r *http.ServeMux) {
	r.Handle(PattenLoglevel, LevelHandler{})
}

func RegisterHandlerWithPath(r *http.ServeMux, path string) {
	r.Handle(path, LevelHandler{})
}

// newDefaultRes admin Default output format.
func newDefaultRes() map[string]interface{} {
	return map[string]interface{}{
		ReturnErrCodeParam: 0,
		ReturnMessageParam: "",
	}
}

// ErrorOutput Unified error output.
func ErrorOutput(w http.ResponseWriter, error string, code int) {
	ret := newDefaultRes()
	ret[ReturnErrCodeParam] = code
	ret[ReturnMessageParam] = error

	_ = json.NewEncoder(w).Encode(ret)
}

// getLevel returns the level of logger's output stream.
func getLevel(logger tlog.Logger, output string) string {
	level := logger.GetLevel(output)
	return tlog.LevelStrings[level]
}

type LevelHandler struct{}

func (h LevelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPut {
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPut}, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorOutput(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		ErrorOutput(w, err.Error(), ErrCodeServer)
		return
	}

	name := r.Form.Get("logger")
	if name == "" {
		name = "default"
	}
	output := r.Form.Get("output")
	if output == "" {
		output = "0" // don't have output, the first output，ordinary users can only configure one.
	}

	logger := tlog.Get(name)
	if logger == nil {
		ErrorOutput(w, "logger not found", ErrCodeServer)
		return
	}

	ret := newDefaultRes()
	if r.Method == http.MethodGet {
		ret["level"] = getLevel(logger, output)
		_ = json.NewEncoder(w).Encode(ret)
	} else if r.Method == http.MethodPut {
		level := r.PostForm.Get("value")

		ret["prelevel"] = getLevel(logger, output)
		logger.SetLevel(output, tlog.LevelNames[level])
		ret["level"] = getLevel(logger, output)

		_ = json.NewEncoder(w).Encode(ret)
	}
}
