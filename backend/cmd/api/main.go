package main

import (
	_ "notes/cmd/api/docs"
	"os"
	"runtime/debug"
)

func main() {
	Init()
	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}
