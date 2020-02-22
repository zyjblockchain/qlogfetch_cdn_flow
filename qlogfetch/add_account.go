package qlogfetch

import (
	"log"
	"os"
	"os/exec"
)

func AddAccount(accessKey, secretKey string) {
	cmd := exec.Command("qlogfetch", "reg", "-ak", accessKey, "-sk", secretKey)
	// 结果打印到命令行
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Print("add_account error: ", err)
		return
	}
	ShowAccount()
}

func ShowAccount() {
	cmd := exec.Command("qlogfetch", "info")
	// 结果打印到命令行
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Print("show_account error: ", err)
		return
	}
}
