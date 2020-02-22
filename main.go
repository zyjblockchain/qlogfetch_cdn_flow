package main

import (
	"flag"
	"fmt"
	"github.com/zyjblockchain/qlogfetch_cdn_flow/handle_logs"
	"github.com/zyjblockchain/qlogfetch_cdn_flow/qlogfetch"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	ak     string // "QTDlNs08wwQQvGLpOgQCssO7csbiPP6hOeCfV0qR"
	sk     string // "YDtKszuj2ymAHLHcJTKOEHNrKyv2zNw2Vp9cqKIz"
	domain string // "video.taped.xiaojing0.com"
	mark   string // "nA0PClSG1LuL_VWTNGJHGkoWm"
	date   string // "2020-02"表示拉取2月份的所有日志 或者 "2020-02-01" 表示只拉取2月1号这一天的日志

	LogsDirPath string // 日志存放路径
)

const (
	reg  = "reg"  // add account
	flow = "flow" // 获取流量

	LogDir = "logs" // 存放下载下来的日志的文件夹
)

func init() {
	// 初始化日志存放路径
	pwd, _ := os.Getwd()
	LogsDirPath = path.Join(pwd, LogDir)
	// 1. 清空logs中的文件
	if err := os.RemoveAll(LogsDirPath); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(LogsDirPath, 0744); err != nil {
		panic(err)
	}
}

func main() {
	// 1. 判断传入的参数
	if len(os.Args) == 1 || (os.Args[1] != reg && os.Args[1] != flow) { // 没有传入任何参数,或者传入的第一个参数不为reg 或者flow
		fmt.Println("参数介绍")
		return
	}
	// 2. 解析传入参数
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&ak, "ak", "", "七牛云上的accessKey")
	flags.StringVar(&sk, "sk", "", "七牛云上的secretKey")
	flags.StringVar(&domain, "domain", "video.taped.xiaojing0.com", "cdn上对应的域名，如: \"video.taped.xiaojing0.com\"")
	flags.StringVar(&mark, "mark", "", "通过特定的标签查询对应的日志信息")
	flags.StringVar(&date, "date", "", "查询日志，\"2020-02\"表示拉取2月份的所有日志 或者 \"2020-02-01\" 表示只拉取2月1号这一天的日志.")
	// 从第二个参数开始解析
	if err := flags.Parse(os.Args[2:]); err != nil {
		log.Print("flag parse error: ", err)
	}
	// 打印出解析出来的参数
	fmt.Println(os.Args[1:])
	fmt.Println(os.Args[1], ak, sk, domain, len(mark), date)

	// 3. 处理不同的业务
	switch os.Args[1] {
	case reg: // add account
		if len(ak) != 0 && len(sk) != 0 {
			qlogfetch.AddAccount(ak, sk)
		} else {
			log.Println("需要输入ak和sk")
		}
		return
	case flow: // 获取某一天的流量
		if date != "" { // "2020-01-11" || "2020-02"
			strs := strings.Split(date, "-")
			// 先解析出年和月
			if len(strs) >= 2 {
				yy, err := strconv.Atoi(strs[0])
				if err != nil {
					log.Println("strconv.Atoi error: ", err)
				}
				mm, err := strconv.Atoi(strs[1])
				if err != nil {
					log.Println("strconv.Atoi error: ", err)
				}

				if len(strs) == 3 {
					dd, err := strconv.Atoi(strs[2])
					if err != nil {
						log.Println("strconv.Atoi error: ", err)
					}
					// 拉取yy-mm-dd的日志
					if err := qlogfetch.DownLogsByDay(yy, mm, dd, domain, LogsDirPath); err != nil {
						log.Println("qlogfetch.DownLogsByDay error: ", err)
						panic(err)
					}
				} else { // ==2 拉取yy-mm的日志
					if err := qlogfetch.DownLogsByMonth(yy, mm, domain, LogsDirPath); err != nil {
						log.Println("qlogfetch.DownLogsByMonth error: ", err)
						panic(err)
					}
				}

				// 提取出流量数量
				flows, err := handle_logs.CalcTotalFlow(LogsDirPath, mark)
				if err != nil {
					panic(err)
				}
				fmt.Printf("流量数量为：%d Byte", flows)
			}
		}
	}

}
