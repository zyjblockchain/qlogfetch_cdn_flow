package handle_logs

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// readAndCalcFlow 读出单个log压缩文件中的内容进行计算
func readAndCalcFlow(logPath string, mark string) (int64, error) {
	flowNum := int64(0)
	// 1. 打开需要解压的文件
	fr, err := os.Open(logPath)
	if err != nil {
		log.Println("open file error: ", err)
		return 0, err
	}

	defer func() {
		if err := fr.Close(); err != nil {
			log.Println("close file error: ", err)
		}
	}()

	// 使用gzip库来读取压缩文件中的内容
	gr, err := gzip.NewReader(fr)
	if err != nil {
		log.Println("error44: ", err)
	}
	defer gr.Close()

	// 使用缓冲流的方式读取大文件
	r := bufio.NewReader(gr)
	for {
		// 读出一行
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Read file by buffer error: ", err)
			return 0, err
		}
		// 提取出这一行中的流量信息
		flowNum += getFlowByStrLine(string(line), mark)
	}

	return flowNum, nil
}

func getFlowByStrLine(line string, mark string) int64 {
	// 判断这行内容中是否存在uuid
	if strings.Contains(line, mark) {
		reg := regexp.MustCompile(`\d\d\d \d+`) // 匹配为前缀是三个数字 + 一个空格 + 后面全是数字的字符串
		result := reg.FindString(line)          // result:  "206 3182157"
		if result == "" {
			return 0
		}
		// 提取流量
		strArr := strings.Split(result, " ")
		if len(strArr) < 2 {
			// 流量字段为空
			return 0
		} else {
			flowStr := strArr[1]
			num, err := strconv.ParseInt(flowStr, 10, 64)
			if err != nil {
				return 0
			}
			return num
		}
	} else {
		return 0
	}
}
