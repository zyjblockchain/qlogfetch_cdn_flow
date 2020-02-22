package qlogfetch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// DownLogsByDay 下载指定的年、月、日和指定cdn域名的日志并保存到指定的路径
func DownLogsByDay(year, month, day int, domain string, savePath string) error {
	mm := strconv.Itoa(month)
	dd := strconv.Itoa(day)
	if month < 10 {
		mm = "0" + mm
	}
	if day < 10 {
		dd = "0" + dd
	}

	date := strconv.Itoa(year) + "-" + mm + "-" + dd
	cmd := exec.Command("qlogfetch", "downlog", "-date", date, "-domains", domain, "-dest", savePath)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}

// DownLogsByMonth 下载指定的当月的所有日志
func DownLogsByMonth(year, month int, domain string, savePath string) error {
	totalDays := 0
	// 判断此月份的天数
	switch month {
	case 2:
		if year%4 == 0 && year%100 != 0 || year%400 == 0 {
			totalDays = 29
		} else {
			totalDays = 28
		}
	case 1, 3, 5, 7, 8, 10, 12:
		totalDays = 31
	case 4, 6, 9, 11:
		totalDays = 30
	default:
		return errors.New(fmt.Sprintf("输入的月份有误: %d", month))
	}

	// 遍历这个月的天数,下载每一天的日志到dirPath路径下
	for d := 1; d <= totalDays; d++ {
		if err := DownLogsByDay(year, month, d, domain, savePath); err != nil {
			return errors.New(fmt.Sprintf("下载日志error: %v", err))
		}
	}
	return nil
}
