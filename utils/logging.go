package utils

import (
	"io/ioutil"
	"log"
	"os"
)

type LogOutput int

const (
	None LogOutput = iota
	StdOut
	File
)

var loggers map[string]*log.Logger = make(map[string]*log.Logger)

func GetLogger(prefix string) *log.Logger {
	logger, exists := loggers[prefix]
	if !exists {
		panic("Logger not found")
		return nil
	}
	return logger
}

func SetupLogger(prefix string) *log.Logger {
	if loggers == nil {
		loggers = make(map[string]*log.Logger)
	}
	logger := log.New(ioutil.Discard, prefix+" ", log.LstdFlags)
	loggers[prefix] = logger
	return logger
}

func SetOutput(prefix string, out LogOutput) {
	logger := GetLogger(prefix)
	if logger == nil {
		return
	}
	switch out {
	case None:
		logger.SetOutput(ioutil.Discard)
	case StdOut:
		logger.SetOutput(os.Stdout)
	case File:
		file, err := os.OpenFile(prefix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err == nil {
			logger.SetOutput(file)
		} else {
			log.Printf("Error cretaing log file for prefix %s: %v", prefix, err)
		}
	}

}

func SetupLoggers() {
	SetupLogger("miner")
	SetupLogger("chain")
	SetupLogger("header")
	SetupLogger("test")
	SetupLogger("main")
	SetupLogger("db")
	//SetOutput("miner", StdOut)
	//SetOutput("header", StdOut)
	//SetOutput("chain", StdOut)
	SetOutput("test", StdOut)
	SetOutput("test", File)
	SetOutput("header", File)
	SetOutput("chain", File)
	SetOutput("db", StdOut)
	SetOutput("main", StdOut)
}
