package message

import (
	"fmt"
	"time"
)

type Message struct {
	To       string
	From     string
	Content  string
	SendTime string
}

// 打印message
func (message *Message) Print() {
	fmt.Printf("%s %s发来消息：\n%s\n", message.SendTime, message.From, message.Content)
}

// 将读到的msg进行分析，设置message的To，Msg和SendTime
func NewMessage(from, to, content string) Message{
	var message = Message{to, from, content, GetNowTime()}
	return message
}

// 获取当前时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}