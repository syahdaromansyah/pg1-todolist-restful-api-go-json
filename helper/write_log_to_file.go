package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func WriteLogToFile(logFn func()) {
	Logger.SetFormatter(&logrus.JSONFormatter{})

	file, _ := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	defer file.Close()

	Logger.SetOutput(file)
	logFn()
}
