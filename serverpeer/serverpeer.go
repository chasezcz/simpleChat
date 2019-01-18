package serverpeer

import (
	"bytes"
	"github.com/multivactech/MultiVAC/rlp"
	"log"
	"net"
	"simpleChat_1.5/message"
	"simpleChat_1.5/peer"
)

type ServerPeer struct {
	Peer peer.Peer
}

// 先对message进行rlp编码，再转交给peer发送字节流。
func (serverPeer *ServerPeer) Send(message message.Message) {
	buf := Encode(message)
	ok := serverPeer.Peer.Send(buf)
	if !ok {
		log.Printf("%s 发送失败", message.To)
	}
}

// 从peer收到字节流，解码成message，返回
func (serverPeer *ServerPeer) Receive() message.Message {
	message := Decode(serverPeer.Peer.Receive())
	return message
}

func (serverPeer *ServerPeer) Close() {
	serverPeer.Peer.Close()
}

// 对message进行rlp编码
func Encode(message message.Message) []byte {
	buf, _ := rlp.EncodeToBytes(message)
	return buf
}

// 对message进行rlp解码
func Decode(buf []byte) message.Message {
	var message message.Message
	err := rlp.Decode(bytes.NewReader(buf), &message)
	if err != nil {
		log.Println(err)
	}
	return message
}

func NewServerPeer(conn net.Conn, ip, port string) ServerPeer {
	return ServerPeer{peer.NewPeer(conn, ip, port)}
}