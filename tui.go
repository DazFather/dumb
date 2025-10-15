package main

import (
	"fmt"

	"github.com/DazFather/brush"
)

const (
	dangerColor  = brush.BrightRed
	warnColor    = brush.BrightYellow
	successColor = brush.BrightGreen
)

type loggable interface {
	Log() string
}

type danger string

func (d danger) Log() string {
	return fmt.Sprintln(brush.Paint(dangerColor, nil, "x"), d)
}

type warn string

func (w warn) Log() string {
	return fmt.Sprintln(brush.Paint(warnColor, nil, "!"), w)
}

func collect(fpath string, logs <-chan loggable) string {
	var msg string

	for log := range logs {
		msg += brush.Paint(brush.BrightMagenta, nil, fpath).String() + " | " + log.Log()
	}

	return msg
}
