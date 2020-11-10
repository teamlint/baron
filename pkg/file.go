package pkg

import "os"

// FileExists 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
