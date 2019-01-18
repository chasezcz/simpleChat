package connmgr

import (
	"log"
	"net"
	"simpleChat_1.5/message"
	"simpleChat_1.5/serverpeer"
	"strings"
	"time"
)

// 保存本地IP和port，拥有一个始终打开的listener，还有一个从ip到peer的映射
type ConnManager struct {
	LocalIP     string
	LocalPort   string
	Listener    net.Listener
	ServerPeers map[string]serverpeer.ServerPeer
}

// 初始化ConnManger，包括本地IP，端口的设置，并开始监听连接
func (connManager *ConnManager) Init(localIP string) bool {
	connManager.LocalIP = localIP
	connManager.LocalPort = "1251" // 可以设置，不一定是1251
	connManager.ServerPeers = make(map[string]serverpeer.ServerPeer)

	addr := connManager.LocalIP + ":" + connManager.LocalPort
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return false
	} else {

		connManager.Listener = listener

		log.Println("开始监听端口")

		return true
	}
}

// 一直监听，等待连接接入，当有连接接入时，得到ip和port，创建一个Peer并存入peers映射中
func (connManager *ConnManager) WaitForConn() {
	for {
		conn, err := connManager.Listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			connManager.NewConnection(conn)
		}
	}
}

// 当有一个正确的conn建立时，对conn进行处理，包装成ServerPeer，并开一个线程开始监听此连接是否有数据发送
func (connManager *ConnManager) NewConnection(conn net.Conn) {
	ip, port := splitAddr(conn.RemoteAddr().String())
	// 封装成ServerPeer，并将之加入到ServerPeers这个映射中
	serverPeer := serverpeer.NewServerPeer(conn, ip, port)

	connManager.ServerPeers[ip] = serverPeer
	// 开一个县城，开始监听其他节点是否发消息
	go ListenToServerPeers(serverPeer)
}

// connManger的关闭函数，包括关闭监听和关闭所有conn
func (connManager *ConnManager) Close() {
	connManager.Listener.Close()
	for _, sp := range connManager.ServerPeers {
		sp.Close()
	}
}

// 循环读取channel中的Message，如果读到，则进入发送环节。
// 发送时，先查map中看是否此ip的连接已建立，若以建立，直接使用，若未建立，先拨号，再发送
func (connManager *ConnManager) ReadyToSend(ch chan message.Message) {
	for {
		message := <- ch
		serverPeer, ok := connManager.ServerPeers[message.To]
		if ok {
			// 因为有可能Send阻塞，所以开一个线程
			go serverPeer.Send(message)
		} else {
			go connManager.DialAndSend(message)
		}
	}
}

// 判断目标地址是否已经建立连接
func (connManager *ConnManager) isExistSP(targetIP string) bool {
	_, ok := connManager.ServerPeers[targetIP]
	return ok
}

// 拨号建立新连接，若连接成功，则发送，否则每隔1s后重试，5次不成功则放弃，报告失败
func (connManager *ConnManager) DialAndSend(message message.Message) {
	ok := false
	for i:= 0; i < 5; i ++ {
		conn, err := net.DialTimeout("tcp", message.To+":1251", time.Second)
		if err != nil {
			time.Sleep(time.Second)
		} else {
			ok = true
			connManager.NewConnection(conn)
			break
		}
	}
	if ok {
		serverPeer := connManager.ServerPeers[message.To]
		serverPeer.Send(message)
	} else {
		log.Printf("%s 连接失败\n", message.To)
	}
}

// 用来循环监听其他人是否发来数据
func ListenToServerPeers(serverPeer serverpeer.ServerPeer) {
	for {
		messages := serverPeer.Receive()
		messages.Print()
	}
}

func splitAddr(addr string) (string, string) {
	return strings.Split(addr, ":")[0], strings.Split(addr, ":")[1]
}