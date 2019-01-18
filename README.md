# 一款简单的局域网聊天工具

## 使用方法
因为此软件使用了1251端口作为listen的固定端口，所以只需要输入想要使用的网络中，本机的ip地址，就开始等待连接了。
之后开始输入命令，命令分两种，普通命令和退出命令
+ 普通命令格式："ip#内容"，IP和内容用#分割
+ 退出命令格式："exit" 即可退出

---

## 设计
这个工具的设计初心是去中心化的，也就是说没有主服务器，采用P2P网络模型，每一个节点都作为tcp的server和client来连接。
设计了ConnManager，ServerPeer，Peer和Message四个主要结构。
+ **ConnManager:** 管理tcp的连接，端口的listen，新连接的发送，断开，以及一些必要的连接逻辑
+ **Peer:** 只负责与节点间字节流的通信，包括发和收
+ **ServerPeer:** 将连接的节点抽象而成，内部含有一个Peer，负责将发送到此节点的Message编码成字节流，或者将从Peer传上来的字节流解析成Message并传到交互界面上
+ **Message：** 发送信息的打包，包括一些必要信息，如from，to，发送内容，发送时间

---

## 原理
https://www.processon.com/view/link/5c41a3f4e4b0fa03ce9f2134
