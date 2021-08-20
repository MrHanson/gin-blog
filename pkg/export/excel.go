package export

import "github.com/MrHanson/gin-blog/pkg/setting"

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}
