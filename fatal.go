// fatal, LitFill <marrazzy54 at email dot com>
// library for fatal assignment or logging (error management)
package fatal

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

func CreateLogFile(filename string) *os.File {
	log_file := Assign(os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	))(
		slog.Default(),
		"cannot create file",
		"file name", filename,
	)
	return log_file
}

func CreateLogger(w io.Writer, lev slog.Level) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(
		io.MultiWriter(w, os.Stderr),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     lev,
		},
	))
	return logger
}

// Debug is a function that wraps slog.
// The log record contains the source position of the caller of Debug.
func Debug(logger *slog.Logger, msg string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Debugf]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
	r.Add(args...)
	err := logger.Handler().Handle(context.Background(), r)
	if err != nil {
		panic(err)
	}
}

// Info is a function that wraps slog.
// The log record contains the source position of the caller of Info.
func Info(logger *slog.Logger, msg string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0])
	r.Add(args...)
	err := logger.Handler().Handle(context.Background(), r)
	if err != nil {
		panic(err)
	}
}

// Error is a function that wraps slog.
// The log record contains the source position of the caller of Error.
func Error(logger *slog.Logger, msg string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	r.Add(args...)
	err := logger.Handler().Handle(context.Background(), r)
	if err != nil {
		panic(err)
	}
}

// Log wraps function call returning error to log it using `log/slog` so it has `msg` and `log`.
// example:
//
//	Log(http.ServeAndListen(port),
//		myLogger,
//	    "Can not serve and listen",
//	    "port", port
//	)
func Log(err error, logger *slog.Logger, msg string, args ...any) {
	if err == nil {
		return
	}
	local_logger := logger.With("err", err)
	Error(local_logger, msg, args...)
	os.Exit(1)
}

// Assign takes the return values of functions that return (val T, err error) and
// returns a function that takes msg and logs in style of slog.logger and return the val T.
//
//	file := Assign(os.Create(fileName))(
//		myLogger,
//	    "cannot create file",
//	    "filename", fileName,
//	)
func Assign[T any](val T, err error) func(*slog.Logger, string, ...any) T {
	if err == nil {
		return func(logger *slog.Logger, msg string, log ...any) T {
			return val
		}
	}
	return func(logger *slog.Logger, msg string, args ...any) T {
		local_logger := logger.With("err", err)
		Error(local_logger, msg, args...)
		return val
	}
}
