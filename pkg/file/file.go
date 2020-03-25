package file

import (
	"os"
	"path"
	"mime/multipart"
	"io/ioutil"
)

// 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// 获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 检查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	// os.IsNotExist：能够接受ErrNotExist、syscall的一些错误，它会返回一个布尔值(true-不存在， false-存在)，能够得知文件不存在或目录不存在
	return os.IsNotExist(err)
}

// 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	// os.IsPermission：能够接受ErrPermission、syscall的一些错误，它会返回一个布尔值(true-无权限， false-有权限)，能够得知权限是否满足
	return os.IsPermission(err)
}

// 如果文件夹不存在，则创建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := Mkdir(src); err != nil {
			return err
		}
	}

	return nil
}

// 创建文件夹
func Mkdir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

