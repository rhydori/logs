package logs_test

import (
	"testing"

	"github.com/rhydori/logs"
)

type discardWriter struct{}

var test = "TESTING"

func (discardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func BenchmarkLog(b *testing.B) {
	logs.SetOutput(discardWriter{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logs.Info("hello world" + test)
	}
}

func BenchmarkLogf(b *testing.B) {
	logs.SetOutput(discardWriter{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logs.Infof("hello world %s", test)
	}
}
