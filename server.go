package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	maplock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server

}

// 监听message广播消息的办法 一旦有消息发送给全部在线用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		//将发送
		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.maplock.Unlock()
	}
}

// 广播消息的办法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//具体业务逻辑
	user := NewUser(conn)

	this.maplock.Lock()
	this.OnlineMap[user.Name] = user
	this.maplock.Unlock()

	this.BroadCast(user, "已上线")

	select {}

}

func (this *Server) Start() {
	//listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Net err:", err)
		return
	}
	//close
	defer listener.Close()

	//启动监听上线消息
	go this.ListenMessager()

	for {
		//accept
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}

		//do habdler
		go this.Handler(con)

	}

}
