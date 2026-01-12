package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type RunFunc func(context.Context) error

type RunError struct {
	Err      error
	ExitCode int
}

func (e RunError) Error() string {
	return e.Err.Error() + fmt.Sprintf(" (exit code %d)", e.ExitCode)
}

// Main is the common executable entry point for all applications.
//
// Any error returned from the provided run function will cause the application
// to exit with the specified exit code. If the error is not of type RunError,
// a default exit code of 1 is used.
func Main(run RunFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx); err != nil {
		rerr, ok := err.(RunError)
		if !ok {
			rerr = RunError{Err: err, ExitCode: 1}
		}

		slog.Error(rerr.Error())
		os.Exit(rerr.ExitCode)
	}
}
