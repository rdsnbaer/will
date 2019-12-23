package loginfo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func init() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath)
	//	fmt.Println(filePath)

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
	go LoggerTimer()
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Println(v)
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

// scheduled tasks, creating log file
func LoggerTimer() {
	ticker := time.NewTicker(time.Minute * 10)
	for _ = range ticker.C {
		tempTime, _ := strconv.Atoi(time.Now().Format("1504"))
		if tempTime >= 0000 && tempTime <= 0020 {
			logger = log.New(openLogFile(getLogFileFullPath()), DefaultPrefix, log.LstdFlags)
		}
	}
}
