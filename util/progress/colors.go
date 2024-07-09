package progress

import (
	"fmt"
	"github.com/fatih/color"
)

type colorFunc func(format string, a ...interface{}) string

var (
	nocolor colorFunc = func(format string, a ...interface{}) string {
		return fmt.Sprintf(format, a)
	}
	DoneColor    colorFunc = color.BlueString
	TimerColor   colorFunc = color.BlueString
	CountColor   colorFunc = color.YellowString
	WarningColor colorFunc = color.BlackString
	SuccessColor colorFunc = color.GreenString
	ErrorColor   colorFunc = color.RedString
	PrefixColor  colorFunc = color.CyanString
)

func NoColor() {
	DoneColor = nocolor
	TimerColor = nocolor
	CountColor = nocolor
	WarningColor = nocolor
	SuccessColor = nocolor
	ErrorColor = nocolor
	PrefixColor = nocolor
}
