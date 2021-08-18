package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

func GetExt(filename string) string {
	return path.Ext(filename)
}

func GetStat(src string) error {
	_, err := os.Stat(src)
	return err
}

func CheckNotExist(src string) bool {
	return os.IsNotExist(GetStat(src))
}

func CheckPremission(src string) bool {
	return os.IsPermission(GetStat(src))
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.Mkdir(src, os.ModePerm)
	if err != nil {
		return err
	}
	return err
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
