package main

import (
	"bufio"         // 用于从连接中读取数据
	"encoding/json" // 用于处理 JSON 编码和解码
	"fmt"           // 用于格式化输出信息
	"net"           // 用于处理网络连接
	"os/exec"       // 用于执行外部命令
	"runtime"       // 用于获取操作系统信息（判断当前系统类型）
)

// Request 结构体，用于解析客户端发送的命令
type Request0 struct {
	Command string `json:"command"` // JSON 中的 "command" 字段将映射到此字段
}

// Response 结构体，用于存储命令执行结果，并发送回客户端
type Response0 struct {
	Result string `json:"result"` // JSON 中的 "result" 字段将映射到此字段
}

// 处理每个客户端连接的函数
func connection(conn net.Conn) {
	defer conn.Close() // 确保连接在函数退出时关闭，避免内存泄漏

	// 打印客户端的远程地址（IP:Port），方便调试
	fmt.Println("客户端已连接：", conn.RemoteAddr())

	// 创建一个读取器，用于从连接中读取客户端发送的数据
	reader := bufio.NewReader(conn)
	for {
		// 从连接中读取客户端发送的 JSON 数据，直到遇到换行符
		jsonData, err := reader.ReadString('\n')
		if err != nil {
			// 如果读取过程中发生错误，打印错误信息并结束循环
			fmt.Println("读取数据失败:", err)
			break
		}

		// 创建一个 Request 结构体，用来接收解析后的 JSON 数据
		var req Request0
		// 解析收到的 JSON 数据，并将其填充到 Request 结构体中
		err = json.Unmarshal([]byte(jsonData), &req)
		if err != nil {
			// 如果 JSON 解析失败，返回错误信息给客户端
			fmt.Println("解析 JSON 失败:", err)
			conn.Write([]byte(`{"result":"Error: Invalid JSON"}` + "\n"))
			continue // 跳过本次循环，继续等待下一条请求
		}

		// 打印收到的命令，便于调试
		fmt.Println("收到命令:", req.Command)

		// 执行客户端发送的命令
		var output []byte
		// 根据操作系统的不同，选择执行命令的方式
		if runtime.GOOS == "windows" {
			// 如果是 Windows 系统，使用 cmd 执行命令
			output, err = exec.Command("cmd", "/C", req.Command).CombinedOutput()
		} else {
			// 如果是类 Unix 系统，使用 sh 执行命令
			output, err = exec.Command("sh", "-c", req.Command).CombinedOutput()
		}

		// 如果命令执行出错，打印错误信息并返回错误内容
		if err != nil {
			fmt.Println("命令执行失败:", err)
			output = []byte("Error: " + err.Error()) // 返回错误信息
		}

		// 创建一个响应结构体，并将命令执行的结果赋值给 Result 字段
		resp := Response0{Result: string(output)}
		// 将响应结构体编码为 JSON 格式
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			// 如果 JSON 编码失败，返回失败信息给客户端
			fmt.Println("生成 JSON 失败:", err)
			conn.Write([]byte(`{"result":"Error: Failed to generate JSON"}` + "\n"))
			continue // 跳过本次循环，继续等待下一条请求
		}

		// 将 JSON 格式的响应数据发送回客户端，并附加换行符表示结束
		_, err = conn.Write(append(jsonResp, '\n'))
		if err != nil {
			// 如果发送数据失败，打印错误信息并退出
			fmt.Println("发送数据失败:", err)
			break
		}
	}
}

// 服务器的入口函数
func main() {
	// 启动一个 TCP 监听器，监听 127.0.0.1:8080 地址
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		// 如果监听失败，打印错误信息并退出
		fmt.Println("监听失败:", err.Error())
		return
	}
	defer listener.Close() // 确保在程序退出时关闭监听器

	// 打印信息，表示服务器已经启动并开始监听指定的端口
	fmt.Println("服务器已启动，正在监听8080端口...")

	// 不断接收客户端的连接
	for {
		// 等待并接受客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			// 如果连接失败，打印错误信息并继续等待下一次连接
			fmt.Println("连接失败:", err.Error())
			continue
		}

		// 每当有新的连接进来时，启动一个新的 goroutine 来处理该连接
		go connection(conn)
	}
}
