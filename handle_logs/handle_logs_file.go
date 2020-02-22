package handle_logs

import (
	"log"
	"path"
	"path/filepath"
)

func CalcTotalFlow(logsDirPath string, mark string) (int64, error) {
	log.Println("log文件路径：", logsDirPath)
	// 获取需要统计的文件列表
	logPathList := getFileListByDir(logsDirPath) // 返回的是文件绝对路径的列表
	log.Println(logPathList)
	totalFlow := int64(0)
	// 遍历日志压缩文件列表
	for _, lp := range logPathList {
		num, err := readAndCalcFlow(lp, mark)
		if err != nil {
			return 0, err
		}
		totalFlow += num
	}
	return totalFlow, nil
}

// getFileListByDir 获取指定文件夹下面的所有文件列表，等价于ls
func getFileListByDir(dirPath string) []string {
	lists, err := filepath.Glob(path.Join(dirPath, "*.gz"))
	if err != nil {
		panic(err)
	}
	return lists
}
