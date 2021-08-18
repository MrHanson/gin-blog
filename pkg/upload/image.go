package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MrHanson/gin-blog/pkg/file"
	"github.com/MrHanson/gin-blog/pkg/logging"
	"github.com/MrHanson/gin-blog/pkg/setting"
	"github.com/MrHanson/gin-blog/pkg/util"
)

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + filepath.Join(GetImagePath(), name)
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		return strings.EqualFold(
			strings.ToUpper(allowExt),
			strings.ToUpper(ext),
		)
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	perm := file.CheckPremission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(filepath.Join(dir, src))
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir src %s, err: %v", src, err)
	}

	return nil
}
