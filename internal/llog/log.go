package llog

import (
	"io"
	"log"
	"os"
	"sync/atomic"
)

const (
	LDEBUG = iota + 1
	LWARN
	LINFO
	LERROR
	LFATAL
)

var logger *Logger = New(os.Stderr, LDEBUG, 0)

type Logger struct {
	level       int64
	w           io.Writer
	debugLogger *log.Logger
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	errLogger   *log.Logger
	fatalLogger *log.Logger
}

func New(w io.Writer, level int64, flag int) *Logger {
	if w == nil {
		w = os.Stderr
	}

	if flag <= 0 {
		flag = log.LstdFlags
	}

	flag = flag | log.Lmsgprefix
	return &Logger{
		w:           w,
		level:       level,
		debugLogger: log.New(w, "[DEBUG] ", flag),
		warnLogger:  log.New(w, "[WARN] ", flag),
		infoLogger:  log.New(w, "[INFO] ", flag),
		errLogger:   log.New(w, "[ERROR] ", flag),
		fatalLogger: log.New(w, "[FATAL] ", flag),
	}
}

func (l *Logger) SetLevel(level int64) {
	if level < LDEBUG || level > LFATAL {
		return
	}

	atomic.StoreInt64(&l.level, level)
}

func (l *Logger) Debugln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LDEBUG {
		return
	}
	l.debugLogger.Println(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LDEBUG {
		return
	}
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LWARN {
		return
	}
	l.warnLogger.Println(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LWARN {
		return
	}
	l.warnLogger.Printf(format, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LINFO {
		return
	}
	l.infoLogger.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LINFO {
		return
	}
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LERROR {
		return
	}
	l.errLogger.Println(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LERROR {
		return
	}
	l.errLogger.Printf(format, v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LFATAL {
		return
	}
	l.fatalLogger.Println(v...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > LFATAL {
		return
	}
	l.fatalLogger.Printf(format, v...)
	os.Exit(1)
}

func Debugln(v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LDEBUG {
		return
	}
	logger.debugLogger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LDEBUG {
		return
	}
	logger.debugLogger.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LWARN {
		return
	}
	logger.warnLogger.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LWARN {
		return
	}
	logger.warnLogger.Printf(format, v...)
}

func Infoln(v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LINFO {
		return
	}
	logger.infoLogger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LINFO {
		return
	}
	logger.infoLogger.Printf(format, v...)
}

func Errorln(v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LERROR {
		return
	}
	logger.errLogger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LERROR {
		return
	}
	logger.errLogger.Printf(format, v...)
}

func Fatalln(v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LFATAL {
		return
	}
	logger.fatalLogger.Println(v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if atomic.LoadInt64(&logger.level) > LFATAL {
		return
	}
	logger.fatalLogger.Printf(format, v...)
	os.Exit(1)
}
