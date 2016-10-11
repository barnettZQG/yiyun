package yiyun

import "os"

//FileIsExit 验证文件是否存在
func FileIsExit(src string) bool {
	_, err := os.Stat(src)
	if err != nil && os.IsExist(err) {
		return false
	}
	return true
}
