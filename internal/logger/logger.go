package logger

import (
	"log"
	"os"

	"github.com/markelca/prioritty/internal/config"
	"github.com/spf13/viper"
)

var logFile *os.File

func InitLogger() error {
	logfile := viper.GetString(config.CONF_LOG_FILE_PATH)
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	logFile = file
	return nil
}

func ShutdownLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
