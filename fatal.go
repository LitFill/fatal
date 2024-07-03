// fatal, LitFill <marrazzy54 at email dot com>
// library for fatal assignment or logging (error management)
package fatal

import (
	"io"
	"log/slog"
	"os"
)

type writer struct {
	files []io.Writer
}

func (w writer) Write(p []byte) (n int, err error) {
	for _, f := range w.files {
		n, err = f.Write(p)
	}
	return
}

func NewWriter(writers []io.Writer) io.Writer {
	return writer{writers}
}

var myWriter = NewWriter([]io.Writer{os.Stderr})

var logger = slog.New(slog.NewJSONHandler(myWriter, &slog.HandlerOptions{AddSource: true}))

func Error(msg string, log ...any) { logger.Error(msg, log...) }
func Debug(msg string, log ...any) { logger.Debug(msg, log...) }
func Info(msg string, log ...any)  { logger.Info(msg, log...) }

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
	Error(err.Error())
	Error(msg, log...)
	os.Exit(1)
}

// Assign takes the return values of functions that return (val T, err error) and
// returns a function that takes msg and logs in style of slog.logger and return the val T.
//
//	file := Assign(os.Create(fileName))(
//	    "cannot create file",
//	    "filename", fileName,
//	)
func Assign[T any](val T, err error) func(string, ...any) T {
	if err == nil {
		return func(msg string, log ...any) T {
			return val
		}
	}
	return func(msg string, log ...any) T {
		Error(err.Error())
		Error(msg, log...)
		return val
	}
}
