package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

// DealResponse 处理server返回消息，显示到标准输出
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		panic(err)
		return
	}
}
func (client *Client) menu() bool {
	var f int
	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	_, err := fmt.Scanln(&f)
	if err != nil {
		panic(err)
		return false
	}
	if f >= 0 && f <= 3 {
		client.flag = f
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字<<<<<")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		panic(err)
		return false
	}
	sendMsg := "rename|" + client.Name + "\n"
	_, err = client.conn.Write([]byte(sendMsg))
	if err != nil {
		return false
	}
	if err != nil {
		fmt.Println("con.Write error:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式...")
		case 2:
			fmt.Println("私聊模式...")
		case 3:
			client.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	// ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}
