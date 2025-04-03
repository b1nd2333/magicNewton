package work

import (
	"fmt"
	"log"
	"magicNewton/api"
	"magicNewton/common"
	"strings"
	"time"
)

func Normal() {
	tokens, err := common.ReadFileLines("token.txt")
	if err != nil {
		fmt.Printf("未找到token.txt文件，请确保文件存在\n")
		return
	}

	proxies, err := common.ReadFileLines("proxy.txt")
	if err != nil {
		fmt.Printf("未找到proxy.txt文件，请确保文件存在。\n")
		return
	}

	proxyCount := len(proxies)

	for {
		for i, v := range tokens {
			proxyStr := proxies[i%proxyCount]
			// 如果ads导出 格式环境编号:token
			api.UserQuests(i+1, v, proxyStr)
		}
		fmt.Println("所有账号领取完毕，等待25小时，再次领取")
		time.Sleep(25 * time.Hour)
	}
}

func ADS() {
	tokens, err := common.ReadFileLines("token.txt")
	if err != nil {
		fmt.Printf("未找到token.txt文件，请确保文件存在\n")
		return
	}

	proxies, err := common.ReadFileLines("proxy.txt")
	if err != nil {
		fmt.Printf("未找到proxy.txt文件，请确保文件存在。\n")
		return
	}

	proxiesMap := make(map[string]string)
	for _, v := range proxies {
		proxy := strings.SplitAfter(v, ":")
		if len(proxy) != 5 {
			fmt.Printf("proxy.txt文件格式不对，格式应为ADS环境编号:ip:port:username:password请检查。\n")
			return
		}
		proxiesMap[proxy[0]] = proxy[1] + ":" + proxy[2] + ":" + proxy[3] + ":" + proxy[4]
	}

	for {
		for i, v := range tokens {
			keys := strings.Split(v, ":")
			if len(keys) == 2 {
				proxyStr := proxiesMap[keys[0]]
				if proxyStr == "" {
					fmt.Printf("账号%d未找到对应编号代理，跳过\n", i+1)
					continue
				}
				api.UserQuests(i+1, keys[1], proxyStr)
			} else {
				log.Fatalln("token.txt格式不对，格式为环境编号:token，或者token，请检查")
			}
		}
		fmt.Println("所有账号领取完毕，等待25小时，再次领取")
		time.Sleep(25 * time.Hour)
	}
}
