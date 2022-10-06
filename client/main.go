package main

import (
	"flag"
	"fmt"
)

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>连接服务器失败")
		return
	}
	//单独开启一个goroutine去处理server返回的消息
	go client.DealResponse()
	fmt.Println(">>>>>连接服务器成功")
	client.Run()
}
