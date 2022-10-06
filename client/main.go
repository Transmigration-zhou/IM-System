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
	fmt.Println(">>>>>连接服务器成功")
	client.Run()
}
