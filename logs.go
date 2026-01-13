package logs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rhydori/logs/ansi"
	"golang.org/x/term"
)

type logEntry struct {
	level []byte
	color []byte
	buf   []byte
}

var out io.Writer = os.Stdout

var (
	logChan = make(chan *logEntry, 1024)
	bufPool = sync.Pool{
		New: func() any {
			b := make([]byte, 0, 256)
			return &b
		},
	}
	reset = []byte("")

	red    = []byte("")
	green  = []byte("")
	yellow = []byte("")
	//blue   = []byte("")
	purple = []byte("")
	cyan   = []byte("")

	lDEBUG = []byte("DEBUG")
	lINFO  = []byte("INFO")
	lWARN  = []byte("WARN")
	lERROR = []byte("ERROR")
	lFATAL = []byte("FATAL")
)

func init() {
	noColorEnv := os.Getenv("NO_COLOR") == "true"
	if !noColorEnv || term.IsTerminal(int(os.Stdout.Fd())) {
		ansi.EnableANSI()
	}

	go writer()
}

func writer() {
	for e := range logChan {
		out.Write(e.buf)
		bufPool.Put(&e.buf)
	}
}

func SetOutput(w io.Writer) {
	out = w
}

func log(level, color []byte, msg string) {
	bp := bufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	now := time.Now()
	buf = append(buf, color...)
	buf = append(buf, '[')
	appendTime(&buf, now)
	buf = append(buf, "] ["...)
	buf = append(buf, level...)
	buf = append(buf, "] "...)
	buf = append(buf, msg...)
	buf = append(buf, reset...)
	buf = append(buf, '\n')

	*bp = buf
	select {
	case logChan <- &logEntry{buf: buf}:
	default:
	}
}

func logf(level, color []byte, format string, args ...any) {
	bp := bufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	now := time.Now()
	buf = append(buf, color...)
	buf = append(buf, '[')
	appendTime(&buf, now)
	buf = append(buf, "] ["...)
	buf = append(buf, level...)
	buf = append(buf, "] "...)

	buf = fmt.Appendf(buf, format, args...)

	buf = append(buf, reset...)
	buf = append(buf, '\n')

	*bp = buf
	select {
	case logChan <- &logEntry{buf: buf}:
	default:
	}
}

func append2(b *[]byte, v int) {
	*b = append(*b, byte('0'+v/10), byte('0'+v%10))
}

func append3(b *[]byte, v int) {
	*b = append(*b,
		byte('0'+v/100),
		byte('0'+(v/10)%10),
		byte('0'+v%10),
	)
}

func appendTime(b *[]byte, t time.Time) {
	h, m, s := t.Hour(), t.Minute(), t.Second()
	ms := t.Nanosecond() / 1e6
	append2(b, h)
	*b = append(*b, ':')
	append2(b, m)
	*b = append(*b, ':')
	append2(b, s)
	*b = append(*b, '.')
	append3(b, ms)
}

func Debug(msg string) { log(lDEBUG, cyan, msg) }
func Info(msg string)  { log(lINFO, green, msg) }
func Warn(msg string)  { log(lWARN, yellow, msg) }
func Error(msg string) { log(lERROR, purple, msg) }
func Fatal(msg string) { log(lFATAL, red, msg); os.Exit(0) }

func Debugf(format string, args ...any) { logf(lDEBUG, cyan, format, args...) }
func Infof(format string, args ...any)  { logf(lINFO, green, format, args...) }
func Warnf(format string, args ...any)  { logf(lWARN, yellow, format, args...) }
func Errorf(format string, args ...any) { logf(lERROR, purple, format, args...) }
func Fatalf(format string, args ...any) { logf(lFATAL, red, format, args...); os.Exit(0) }
