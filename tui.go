package main

import (
	"fmt"

	"github.com/DazFather/brush"
)

var (
	red     = brush.New(brush.BrightRed, nil)
	yellow  = brush.New(brush.BrightYellow, nil)
	magenta = brush.New(brush.BrightMagenta, nil)
)

func caret(line string, at int) string {
	return fmt.Sprintf("%s\n%*s%s\n", line, at-1, "", red.Paint("^"))
}

func danger(v ...any) string {
	return red.Paint(" x ").String() + fmt.Sprintln(v...)
}

func warn(v ...any) string {
	return yellow.Paint(" ! ").String() + fmt.Sprintln(v...)
}

func collect(fpath string, logs ...string) string {
	var msg string

	for _, log := range logs {
		msg += magenta.Paint(fpath).String() + " |" + log
	}

	return msg + "\n"
}
