package peer

import (
	"log"
	"net"
)

type Peer struct {
	Conn 	   net.Conn
	TargetIP   string
	TargetPort string
}

// 从Peer方收到信息，返回字节流
func (peer *Peer) Receive() []byte{
	buf := make([]byte, 1024)
	n, err := peer.Conn.Read(buf)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return buf[:n]
	}
}

// 发送给peer字节流，如果发送不成功，等待一秒后继续尝试，5次后则报告发送失败
func (peer *Peer) Send(buf []byte) bool{
	ok := false
	for i := 0; i < 5; i ++{
		_, err := peer.Conn.Write(buf)
		if err == nil {
			ok = true
			break
		}
	}
	return ok
}

// 关闭Peer和连接
func (peer *Peer) Close() {
	peer.Conn.Close()
}

// 新建一个Peer，返回Peer
func NewPeer(conn net.Conn, ip,  port string)  Peer{
	return Peer{conn, ip, port}
}
