package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strings"
	"time"
)

type RespStruct struct {
	Data []struct {
		Id            string    `json:"id"`
		UserId        string    `json:"userId"`
		QuestId       string    `json:"questId"`
		Status        string    `json:"status"`
		Credits       int       `json:"credits"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
		DiceRolls     []int     `json:"_diceRolls,omitempty"`
		BonusFactor   int       `json:"_bonus_factor,omitempty"`
		RolledCredits int       `json:"_rolled_credits,omitempty"`
	} `json:"data"`
}

// UserQuests 获取上次完成时间，是否24小时
func UserQuests(num int, token, proxyStr string) error {
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 创建请求
	req, err := http.NewRequest("GET", "https://www.magicnewton.com/portal/api/userQuests", nil)
	if err != nil {
		fmt.Println(err)
		return err
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
			if strings.Contains(string(body), "Invalid session") {
				color.Red("账号%dSession已过期\n", num)
				return nil
			}
		}
		return errors.New(resp.Status)
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return err
	}

	// 检查是否是 Gzip 压缩
	if resp.Header.Get("Content-Encoding") == "br" {
		body, err = decompressBrotli(body)
		if err != nil {
			fmt.Println("解压 Gzip 失败:", err)
			return err
		}
	}

	respModel := &RespStruct{}
	json.Unmarshal(body, respModel)

	//fmt.Println(string(body))
	currentUTC := time.Now().UTC()
	var lastCheckTime time.Time
	//roll(num, token, proxyStr)
	for _, v := range respModel.Data {
		if v.QuestId == "f56c760b-2186-40cb-9cbc-3af4a3dc20e2" { //最后一次签到时间
			lastCheckTime = v.CreatedAt.Add(24 * time.Hour)
		}
	}

	//lastCheckTime := respModel.Data[len(respModel.Data)-1].CreatedAt.Add(24 * time.Hour)
	if currentUTC.After(lastCheckTime) { // 在之后，可以签到
		err = roll(num, token, proxyStr)
		if err != nil {
			return err
		}
	} else { // 不能签到
		color.Yellow("账号%d还未到达签到时间，下次时间为UTC时间%s\n", num, lastCheckTime.Format(time.RFC3339))
	}

	return nil
}

func parseProxy(account string) (ip, port, username, password string) {
	// 假设 proxy 格式为 "ip:port:username:password"
	parts := strings.Split(account, ":")
	ip, port = parts[0], parts[1]
	if len(parts) > 2 {
		username = parts[2]
	}
	if len(parts) > 3 {
		password = parts[3]
	}

	return
}

// 对 Brotli 压缩的数据进行解压缩
func decompressBrotli(compressedData []byte) ([]byte, error) {
	// 创建一个 bytes.Reader 用于读取压缩数据
	reader := bytes.NewReader(compressedData)
	// 创建 Brotli 解码器
	brotliReader := brotli.NewReader(reader)
	// 创建一个 bytes.Buffer 用于存储解压缩后的数据
	var decompressedBuffer bytes.Buffer
	// 将解压缩后的数据写入 bytes.Buffer
	_, err := io.Copy(&decompressedBuffer, brotliReader)
	if err != nil {
		return nil, err
	}
	// 返回解压缩后的数据
	return decompressedBuffer.Bytes(), nil
}
