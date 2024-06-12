// fatal, LitFill <marrazzy54 at email dot com>
// library for fatal assignment or logging (error management)
package fatal

import (
	"log/slog"
	"os"
)

// the logger used by package level functions
var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}))

// Log wraps function call returning error to log it using `log/slog` so it has `msg` and `log`.
// example:
//
//	Log(http.ServeAndListen(port),
//	    "Can not serve and listen",
//	    "port", port
//	)
func Log(err error, msg string, log ...any) {
	if err == nil {
		return
	}
	logger.Error(err.Error())
	logger.Error(msg, log...)
	os.Exit(1)
}

// Assign wraps function call returning a `val` and error to log it
// if error != nil using `log/slog` so it has `msg` and `log`.
// example:
//
//	file := Assign(os.Create("log.txt"),
//	    "Can not create file",
//	    "file", "log.txt"
//	)
func Assign[T any](val T, err error, msg string, log ...any) T {
	if err == nil {
		return val
	}
	logger.Error(err.Error())
	logger.Error(msg, log...)
	return val
}
