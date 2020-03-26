package logging


import (
	"os"
	"time"
	"fmt"

	"github.com/Songkun007/go-gin-blog/pkg/setting"
	"github.com/Songkun007/go-gin-blog/pkg/file"
)

// 获取日志存储路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)	// %s - 直接输出原始字符串
}

// 获取日志名
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

// 获取文件句柄
func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}


	/* os.OpenFile：调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。
	   如果出现错误，则为*PathError

	   const (
	    // Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	    O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
	    O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
	    O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
	    // The remaining values may be or'ed in to control behavior.
	    O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
	    O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
	    O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
	    O_SYNC   int = syscall.O_SYNC   // 同步IO
	    O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时
	)
	   0644 代表文件的模式和权限位，r(4)、w(2)、x(1)
	*/
	f, err := file.Open(src + fileName, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
