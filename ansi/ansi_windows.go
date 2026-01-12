//go:build windows

package ansi

import (
	"os"

	"golang.org/x/sys/windows"
)

func EnableANSI() {
	h := windows.Handle(os.Stdout.Fd())

	var mode uint32
	windows.GetConsoleMode(h, &mode)
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	windows.SetConsoleMode(h, mode)
}
