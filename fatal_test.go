package fatal

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

func TestAssign(t *testing.T) {
	type args[T any] struct {
		val T
		err error
		msg string
	}
	tests := []struct {
		args args[any]
		want any
		name string
	}{
		{
			name: "12",
			args: args[any]{12, nil, "msg"},
			want: 12,
		},
		{
			name: "err",
			args: args[any]{0, errors.New("testing error case"), "msg"},
			want: 0,
		},
		{
			name: "struct",
			args: args[any]{struct{ num int }{13}, nil, "msg"},
			want: struct{ num int }{13},
		},
		{
			name: "slice",
			args: args[any]{[]int{1, 2, 3}, nil, "msg"},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Assign(tt.args.val, tt.args.err)(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assign()() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestLog(t *testing.T) {
	type args struct {
		err error
		msg string
		log []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil",
			args: args{err: nil, msg: "msg"},
		},
		{
			name: "msg with nil",
			args: args{err: nil, msg: "new custom message"},
		},
		{
			name: "nil with msg and log",
			args: args{err: nil, msg: "another new message", log: []any{"number", 12}},
		},
		// {
		// 	name: "not nil err",
		// 	args: args{err: errors.New("custom error"), msg: "another new with err message"},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.err, tt.args.msg, tt.args.log...)
		})
	}
}

func Test_logger_Info(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		l    *slog.Logger
		args args
	}{
		{
			name: "logger.Info",
			l:    slog.New(slog.NewJSONHandler(myWriter, nil)),
			args: args{"logger.Info(msg)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Info(tt.args.msg)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normael",
			args: args{"info msg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.msg)
		})
	}
}

func Test_logger_Debug(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		l    *slog.Logger
		args args
	}{
		{
			name: "logger.Debug",
			l:    slog.New(slog.NewJSONHandler(myWriter, nil)),
			args: args{"logger.Debug(msg)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debug(tt.args.msg)
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{"debug msg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.msg)
		})
	}
}
