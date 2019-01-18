package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"simpleChat_1.5/connmgr"
	"simpleChat_1.5/message"
	"strconv"
	"strings"
)

func main() {
	connManager := new(connmgr.ConnManager)

	// 重输本地ip直到成功
	localIp := getLocalIp()
	for !connManager.Init(localIp) {
		localIp = getLocalIp()
	}
	// 开一个线程让其一直响应连接
	go connManager.WaitForConn()
	// 主线程用来交互
	Interaction(connManager)

	defer connManager.Close()
}

// 交互函数，包括监听键盘是否有输入（回车发送），并使用channel作为线程间通讯的方式
func Interaction(connManager *connmgr.ConnManager) {
	ch := make(chan message.Message)
	go connManager.ReadyToSend(ch)

	for {
		command := scanToLine()
		// 符合发送命令规范
		if len(command) > 0 && strings.Contains(command, "#") {
			ip, content := splitCommand(command)

			if CheckAddr(ip) {
				// 交给Message打包
				messages := message.NewMessage(connManager.LocalIP, ip, content)
				// 传给channel，由ReadToSend进程发送
				ch <- messages
			} else {
				log.Println("输入有误")
			}
		} else if command == "exit" {  // 退出命令
			log.Println("再见")
			break
			connManager.Close()
		} else {
			log.Println("输入有误")
		}
	}
}

// 将普通命令拆分成 目标ip 和 要发送的内容
func splitCommand(command string) (string, string) {
	ss := strings.Split(command, "#")
	return ss[0], ss[1]
}

// 获得本机ip
func getLocalIp() string {
	fmt.Println("请输入本机ip地址：")
	return scanToLine()
	//return "192.168.12.103"
}

// 每次读一行，返回string
func scanToLine() string{
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}

// 检查地址是否有误，包括ip （和port）
func CheckAddr(addr string) bool{
	ip := strings.Split(addr, ":")[0]
	//port := strings.Split(addr, "：")[1]
	for _,valueString := range strings.Split(ip, ".") {
		value, err := strconv.Atoi(valueString)
		if err != nil {
			log.Println(err)
			return false
		} else if value < 0 || value > 254{
			return false
		}
	}
	//value, err := strconv.Atoi(port)
	//if err != nil {
	//	log.Fatal(err)
	//} else if value < 0 || value > 65535 {
	//	return false
	//}
	return true
}