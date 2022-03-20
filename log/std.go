package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"

	klog "github.com/go-kratos/kratos/v2/log"
)

var _ klog.Logger = (*stdLogger)(nil)

// DefaultLogger is default logger.
var DefaultStdLogger klog.Logger = NewStdLogger(log.Writer())

func GetLogger() klog.Logger {
	return DefaultStdLogger
}

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) klog.Logger {
	return &stdLogger{
		log: log.New(w, "", log.LstdFlags|log.Llongfile),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level klog.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	_ = l.log.Output(4, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}
