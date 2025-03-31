package main

import (
	"bufio"
	"fmt"
	"magicNewton/api"
	"os"
	"time"
)

func main() {
	tokens, err := readFileLines("token.txt")
	if err != nil {
		fmt.Printf("未找到token.txt文件，请确保文件存在\n")
		return
	}

	proxies, err := readFileLines("proxy.txt")
	if err != nil {
		fmt.Printf("未找到proxy.txt文件，请确保文件存在。\n")
		return
	}

	proxyCount := len(proxies)
	for {
		for i, v := range tokens {
			proxyStr := proxies[i%proxyCount]
			api.UserQuests(i+1, v, proxyStr)
		}
		fmt.Println("所有账号领取完毕，等待25小时，再次领取")
		time.Sleep(25 * time.Hour)
	}

}

// 读取文件内容，返回每行内容的切片
func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
