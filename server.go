package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	Maplock   sync.RWMutex

	message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,

		OnlineMap: make(map[string]*User),
		message:   make(chan string),
	}
	return server
}

func (this *Server) BoardCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection Successed!" + conn.RemoteAddr().String())

	user := NewUser(conn, this)

	user.Online()

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("err:", err)
				return

			}
			msg := string(buf[:n-1])
			fmt.Println("[" + user.Addr + "]" + user.Name + ":" + msg)
			user.DoMessage(msg)

		}

	}()

	select {}
}

func (this *Server) ListenMessager() {
	for {
		sendMsg := <-this.message

		this.Maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- sendMsg
		}
		this.Maplock.Unlock()

	}

}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	defer listener.Close()

	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err", err)
			continue
		}

		go this.Handler(conn)
	}

}
