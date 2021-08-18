package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	F, err := openLogFile(getLogFilePath(), getLogFileName())
	if err != nil {
		log.Fatalln(err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func logPrint(args ...interface{}) {
	for _, v := range args {
		logger.Println(v)
	}
}

func Debug(args ...interface{}) {
	setPrefix(DEBUG)
	logPrint(args)
}

func Info(args ...interface{}) {
	setPrefix(INFO)
	logPrint(args)
}

func Warn(args ...interface{}) {
	setPrefix(WARNING)
	logPrint(args)
}

func Error(args ...interface{}) {
	setPrefix(ERROR)
	logPrint(args)
}

func Fatal(args ...interface{}) {
	setPrefix(FATAL)
	logPrint(args)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
