// fatal, LitFill <marrazzy54 at gmail dot com>
// library for fatal assignment or logging (error management)
package fatal

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

// CreateLogFile creates a new log file with the specified filename and
// returns a pointer to the opened file.
// It uses Assign with the default slog.Logger to handle the assignment.
func CreateLogFile(filename string) *os.File {
	return Assign(os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	))(
		slog.Default(),
		"cannot create file",
		"file name", filename,
	)
}

// CreateLogger creates a new logger object using the slog package and
// configures it to write logs to the provided writer `w`
// with the specified logging level `lev`.
func CreateLogger(w io.Writer, lev slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     lev,
		},
	))
}

// Debug logs an error message with the source position of the caller.
func Debug(logger *slog.Logger, msg string, args ...any) {
	logMessage(context.Background(), logger, slog.LevelDebug, msg, args...)
}

// Info logs an error message with the source position of the caller.
func Info(logger *slog.Logger, msg string, args ...any) {
	logMessage(context.Background(), logger, slog.LevelInfo, msg, args...)
}

// Error logs an error message with the source position of the caller.
func Error(logger *slog.Logger, msg string, args ...any) {
	logMessage(context.Background(), logger, slog.LevelError, msg, args...)
}

// logMessage is a helper function that logs a message with given level and
// source position of the caller.
func logMessage(
	ctx context.Context,
	logger *slog.Logger, level slog.Level,
	msg string, args ...any,
) {
	if !logger.Enabled(ctx, level) {
		return
	}
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc) // 3 is runtime.Callers(), logMessage(), and Error()/Info/Debug
	pc = pc[:n]
	start, _ := filterPc(pc)
	r := slog.NewRecord(time.Now(), level, msg, pc[start])
	r.Add(args...)
	if err := logger.Handler().Handle(ctx, r); err != nil {
		panic(err)
	}
}

// filterPc filters the pc slice for frames that are not in this package
func filterPc(pc []uintptr) (start, end int) {
	if len(pc) == 0 {
		return
	}
	frames := runtime.CallersFrames(pc)
	indexes := make([]int, 0)
	counter := 0
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		if !strings.Contains(frame.Function, "fatal") {
			indexes = append(indexes, counter)
		}
		counter++
	}
	return indexes[0], indexes[len(indexes)-1]
}

// Log logs the provided error with the specified message and logger, then exits the program.
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
	localLogger := logger.With("err", err)
	Error(localLogger, msg, args...)
	os.Exit(1)
}

// Assign handles the return values of functions that return (val T, err error)
// and returns a function that logs the error and message using slog.Logger and
// returns val T. example:
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
		localLogger := logger.With("err", err)
		Error(localLogger, msg, args...)
		return val
	}
}
