package main

func main() {
	server := NewServer("127.0.1.1", 8888)
	server.Start()
}
