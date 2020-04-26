package export

import "github.com/Songkun007/go-gin-blog/pkg/setting"

// 获取文件访问url路径
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}


func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// 获取存储全路径
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
