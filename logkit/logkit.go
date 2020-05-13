package logkit

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	info    *log.Logger
	warning *log.Logger
	erring  *log.Logger
)

func Info(data ...interface{}) {
	info.Println(data)
}
func InfoF(format string, data ...interface{}) {
	info.Printf(format, data)
}

func Err(data ...interface{}) {
	erring.Println(data)
}
func ErrF(format string, data ...interface{}) {
	erring.Printf(format, data)
}
func Warn(data ...interface{}) {
	warning.Println(data)
}

func Warnf(format string, data ...interface{}) {
	warning.Printf(format, data)
}
func LogInit(logdir string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if logdir[:2] == "./" {
		dir += logdir[1:]
	}
	errFile, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.LUTC)
	warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	erring = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}
