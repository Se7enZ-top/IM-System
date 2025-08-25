package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user

}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.Maplock.Lock()

	this.server.OnlineMap[this.Name] = this

	this.server.Maplock.Unlock()

	this.server.BoardCast(this, "已上线")

}

func (this *User) Offline() {
	this.server.Maplock.Lock()

	delete(this.server.OnlineMap, this.Name)

	this.server.Maplock.Unlock()

	this.server.BoardCast(this, "下线")

}

func (this *User) DoMessage(msg string) {
	this.server.BoardCast(this, msg)
}
