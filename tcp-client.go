package main

import (
	"bufio"         // 用于读取用户输入和从网络连接中读取数据
	"encoding/json" // 用于处理 JSON 编码和解码
	"fmt"           // 用于格式化输出
	"net"           // 用于网络通信，特别是客户端和服务器之间的 TCP 连接
	"os"            // 用于处理操作系统相关的操作，这里用于获取用户输入
)

// Request1 结构体，用于构造客户端发送的请求数据
type Request1 struct {
	Command string `json:"command"` // "command" 字段映射到结构体中的 Command 字段
}

// Response1 结构体，用于解析服务端返回的响应数据
type Response1 struct {
	Result string `json:"result"` // "result" 字段映射到结构体中的 Result 字段
}

func main() {
	// 连接到本地的 8080 端口
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		// 如果连接失败，输出错误并退出
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close() // 确保在函数退出时关闭连接

	// 输出连接成功的消息
	fmt.Println("已连接到服务器")

	// 创建一个读取器，用于从标准输入读取数据
	reader := bufio.NewReader(os.Stdin)
	for {
		// 提示用户输入命令
		fmt.Print("请输入命令：")
		command, _ := reader.ReadString('\n') // 读取用户输入的命令
		command = command[:len(command)-1]    // 去除命令末尾的换行符

		// 构造 JSON 请求
		req := Request1{Command: command}
		// 将请求结构体转化为 JSON 格式
		jsonReq, err := json.Marshal(req)
		if err != nil {
			// 如果生成 JSON 失败，输出错误并继续下一次循环
			fmt.Println("生成 JSON 失败:", err)
			continue
		}

		// 发送 JSON 请求到服务器
		_, err = conn.Write(append(jsonReq, '\n')) // 向服务器发送请求，并附加换行符
		if err != nil {
			// 如果发送失败，输出错误并退出
			fmt.Println("发送命令失败:", err)
			return
		}

		// 读取服务器返回的 JSON 响应
		jsonResp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// 如果接收响应失败，输出错误并退出
			fmt.Println("接收响应失败:", err)
			return
		}

		// 解析服务器返回的 JSON 响应
		var resp Response1
		err = json.Unmarshal([]byte(jsonResp), &resp)
		if err != nil {
			// 如果解析 JSON 失败，输出错误并继续下一次循环
			fmt.Println("解析 JSON 失败:", err)
			continue
		}

		// 打印服务器的响应结果
		fmt.Println("服务器响应:\n", resp.Result)
	}
}
