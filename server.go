package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// ListenMessage 监听Message广播消息channel的goroutine，有消息就直接发送给全部的在线user
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// BroadCast 广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
	user := NewUser(conn, s)
	user.Online()
	//监听用户是否活跃的channel
	isLive := make(chan bool)
	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("con read Error", err)
				return
			}
			//提取用户消息（去除'\n'）
			msg := string(buf[:n-1])
			//用户针对msg进行消息处理
			user.DoMessage(msg)
			//用户的任意消息代表当前用户是活跃的
			isLive <- true
		}
	}()
	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃，重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 10):
			// 已经过时
			// 将当前用户强制删除
			user.SendMsg("超过10秒未操作，你被强制下线")
			// 销毁资源
			close(user.C)
			// 退出当前handler
			return
		}
	}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(listener)
	//启动监听Message的goroutine
	go s.ListenMessage()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
