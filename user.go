package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

// Online 用户上线的业务
func (u *User) Online() {
	//用户上线，加入到OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// Offline 用户上线的业务
func (u *User) Offline() {
	//用户下线，从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播当前用户上线消息
	u.server.BroadCast(u, "下线")
}

// SendMsg 给当前用户对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}

// DoMessage 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式:rename|xx
		newName := msg[7:]
		if newName == "" {
			u.SendMsg("用户名不能为空\n")
			return
		}
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg(newName + "已被占用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("修改成功，当前用户名:" + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式:to|xx|消息内容
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("用户名不能为空\n")
			return
		}
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg(remoteName + "不存在\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("发送消息内容不能为空\n")
			return
		}
		remoteUser.SendMsg(u.Name + "对你说:" + content + "\n")
	} else {
		u.server.BroadCast(u, msg)
	}
}

// ListenMessage 监听当前User channel的方法，有消息就直接发送给客户端
func (u *User) ListenMessage() {
	for msg := range u.C {
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			panic(err)
		}
	}
	err := u.conn.Close()
	if err != nil {
		panic(err)
	}
}
