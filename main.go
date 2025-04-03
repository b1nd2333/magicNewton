package main

import (
	"bufio"
	"fmt"
	"magicNewton/work"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 只运行一次选择逻辑
	fmt.Println("============================")
	fmt.Println("🎯 请选择模式：")
	fmt.Println("1️⃣  普通模式")
	fmt.Println("2️⃣  ADS导出模式")
	fmt.Println("============================")
	fmt.Print("👉 请输入序号（1 / 2）：")

	// 读取用户输入
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // 去除空格和换行符

	// 选择对应模式
	if input == "1" {
		work.Normal()
	} else if input == "2" {
		work.ADS()
	} else {
		fmt.Println("❌ 输入错误，请重新运行程序并选择正确的模式！")
		os.Exit(1)
	}

}
