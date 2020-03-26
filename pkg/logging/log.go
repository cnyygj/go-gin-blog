package logging

import (
	"os"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix = ""
	DefaultCallerDepth = 2

	logger *log.Logger
	logPrefix = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota	// 每次 const 出现时，都会让 iota 初始化为0，后面的参数会自动加一
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath()						// 获取日志存储路径
	fileName := getLogFileName()						// 获取日志名
	F, err = openLogFile(fileName, filePath)			// 获取文件句柄
	if err != nil {
		log.Fatalln(err)
	}

	// 创建一个logger 参数1：日志写入目的地， 参数2：每条日志的前缀， 参数3：日志属性
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// ... 不确定数量传参
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

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

// 设置日志输出格式
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth) // 获取文件信息以及函数运行行数
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}