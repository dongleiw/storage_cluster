package common

import (
	"fmt"
	builtinlog "log"
	"path"
	"runtime"
	"sync"
)

var (
	DEBUG  = 1
	INFO   = 2
	NOTICE = 3
	WARN   = 4
	ERROR  = 5
	FATAL  = 6
)

func SetLogLevel_str(log_level string) error {
	switch log_level {
	case "debug", "DEBUG":
		g_log_wrapper.log_level = DEBUG
	case "info", "INFO":
		g_log_wrapper.log_level = INFO
	case "notice", "NOTICE":
		g_log_wrapper.log_level = NOTICE
	case "warn", "WARN":
		g_log_wrapper.log_level = WARN
	case "error", "ERROR":
		g_log_wrapper.log_level = ERROR
	case "fatal", "FATAL":
		g_log_wrapper.log_level = FATAL
	default:
		return fmt.Errorf("unknown log_level[%v]", log_level)
	}
	return nil
}

func Log_init() {
	builtinlog.SetFlags(builtinlog.Ldate | builtinlog.Ltime | builtinlog.Lmicroseconds)
}

func Debug(fmts string, v ...interface{})  { log(DEBUG, fmts, 2, v...) }
func Info(fmts string, v ...interface{})   { log(INFO, fmts, 2, v...) }
func Notice(fmts string, v ...interface{}) { log(NOTICE, fmts, 2, v...) }
func Warn(fmts string, v ...interface{})   { log(WARN, fmts, 2, v...) }
func Error(fmts string, v ...interface{})  { log(ERROR, fmts, 2, v...) }
func Fatal(fmts string, v ...interface{})  { log(FATAL, fmts, 2, v...) }

func Debug2(fmts string, v ...interface{})  { log(DEBUG, fmts, 3, v...) }
func Info2(fmts string, v ...interface{})   { log(INFO, fmts, 3, v...) }
func Notice2(fmts string, v ...interface{}) { log(NOTICE, fmts, 3, v...) }
func Warn2(fmts string, v ...interface{})   { log(WARN, fmts, 3, v...) }
func Error2(fmts string, v ...interface{})  { log(ERROR, fmts, 3, v...) }
func Fatal2(fmts string, v ...interface{})  { log(FATAL, fmts, 3, v...) }

type _LogWrapper struct {
	log_level int
	lock      sync.Mutex
}

var (
	level_str_array = []string{
		"",
		"[DEBUG]  ",
		"[INFO]   ",
		"[NOTICE] ",
		"[WARN]  ",
		"[ERROR]  ",
		"[FATAL]  ",
	}

	g_log_wrapper = _LogWrapper{
		log_level: DEBUG,
	}
)

func log(level int, fmts string, call_level int, v ...interface{}) {
	if level < g_log_wrapper.log_level {
		return
	}

	var pc, file, line, ok = runtime.Caller(call_level)
	var fnname = "FNNAME"

	if ok {
		file = path.Base(file)
		var f = runtime.FuncForPC(pc)
		if f != nil {
			fnname = f.Name()
		}
	}

	g_log_wrapper.lock.Lock()
	defer g_log_wrapper.lock.Unlock()

	builtinlog.SetPrefix(level_str_array[level])
	builtinlog.Printf("[%v:%v] [%v] %v", file, line, fnname, fmt.Sprintf(fmts, v...))
}
