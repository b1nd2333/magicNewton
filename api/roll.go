package api

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strings"
	"time"
)

type rollRespStruct struct {
	Message string `json:"message"`
	Data    struct {
		Id          string    `json:"id"`
		UserId      string    `json:"userId"`
		QuestId     string    `json:"questId"`
		Status      string    `json:"status"`
		Credits     int       `json:"credits"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		BonusFactor int       `json:"_bonus_factor"`
	} `json:"data"`
}

func roll(num int, token, proxyStr string) error {
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return roll(num, token, proxyStr)
	}

	// 创建请求
	bodyStr := "{\n  \"questId\": \"f56c760b-2186-40cb-9cbc-3af4a3dc20e2\",\n  \"metadata\": {}\n}"
	req, err := http.NewRequest("POST", "https://www.magicnewton.com/portal/api/userQuests", strings.NewReader(bodyStr))
	if err != nil {
		fmt.Println(err)
		return roll(num, token, proxyStr)
	}

	req.Header = createHeaders(token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			// 读取响应数据
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("读取响应体出错: %v\n", err)
				return err
			}

			if strings.Contains(string(body), "Quest already completed") {
				fmt.Printf("账号%d今日已签到\n", num)
				return nil
			}
		}

		return err
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return roll(num, token, proxyStr)
	}

	// 检查是否是 Gzip 压缩
	if resp.Header.Get("Content-Encoding") == "br" {
		body, err = decompressBrotli(body)
		if err != nil {
			fmt.Println("解压 Gzip 失败:", err)
			return roll(num, token, proxyStr)
		}
	}

	rollRespModel := &rollRespStruct{}
	err = json.Unmarshal(body, rollRespModel)
	if rollRespModel.Message == "Quest completed" {
		color.Green("账号%d投掷成功，获得%d积分\n", num, rollRespModel.Data.Credits)
	} else {
		color.Red("账号%d投掷失败，%s\n", num, rollRespModel.Message)
	}
	return nil
}
