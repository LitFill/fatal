// fatal, LitFill <marrazzy54 at email dot com>
// library for fatal assignment or logging (error management)
package fatal

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/OlegStotsky/go-monads/either"
)

type Reporter interface {
	Report(msg string, log ...any)
}

// the logger used by package level functions
var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}))

type reporter struct{ logger *slog.Logger }

func (r reporter) Report(msg string, log ...any) { r.logger.Error(msg, log...) }
func newReporter(logger *slog.Logger) reporter   { return reporter{logger} }

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

func LogV2(err error) Reporter {
	if err == nil {
		return nil
	}
	return newReporter(logger)
}

type myResult[T any] struct {
	Ok    T
	Error error
}

func NewResult[T any](ok T, error error) myResult[T] {
	return myResult[T]{
		Ok:    ok,
		Error: error,
	}
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

func AssignWithEither[T any](val either.Either[error, T], msg string, log ...any) T {
	value, err := either.ToMaybe[error, T](val).Get()
	if err == nil {
		return value
	}
	logger.Error(err.Error())
	logger.Error(msg, log...)
	return value
}

func AssignWithResult[T any](result myResult[T], msg string, log ...any) T {
	if result.Error == nil {
		return result.Ok
	}
	logger.Error(result.Error.Error())
	logger.Error(msg, log...)
	return result.Ok
}

func AssignV2[T any](val T, err error, msg string, log ...any) T {
	if err == nil {
		return val
	}
	logger.Error(err.Error())
	logger.Error(msg, log...)
	return val
}

func wrap[T any](val T, err error) T {
	if err == nil {
		return val
	}
	os.Exit(1)
	return val
}

func myErr() error { return fmt.Errorf("my own error") }
func getInt() (int, error) {
	return 0, nil
}

type rrr[T any] interface {
	Result[T, error]
}

func Errorable[T any](val T, err error) rrr[T] {
	return FromTuple(val, err)
}

func AssignR[T any](r rrr[T], msg string, log ...any) T {
	if r.Failure() {
		logger.Error(r.Error().Error())
		logger.Error(msg, log...)
	}
	return r.Value()
}

func main() {
	assign := AssignWithResult[int]
	n := NewResult[int]

	great := assign(n(getInt()), "msg", "k", 0)
	fmt.Println(great)

	AssignWithResult(NewResult(os.Create("log.json")), "Message", "file", "log.json")
	AssignR(Errorable(os.Create("file")), "", "", 0)
	LogV2(myErr()).Report("msg", "log", "logged")

	m := AssignWithEither[int](either.FromErrorable(getInt()), "msg", "k", 0)
	fmt.Println(great + m)

	i := AssignR(Errorable(getInt()), "", "", 0)
	fmt.Println(i)
}
