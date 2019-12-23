package loginfo

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "200601"
	TimeFormats = "20060102"
)

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormats), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath+"/", suffixPath)
}

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath+time.Now().Format(TimeFormat))
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission: %v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatal("Fail to Openfile :%v", err)
	}

	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
