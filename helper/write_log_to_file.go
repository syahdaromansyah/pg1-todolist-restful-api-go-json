package helper

import (
	"log"
	"os"
)

func WriteLogToFile(logFn func()) {
	file, _ := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	defer file.Close()

	log.SetOutput(file)
	logFn()
}
