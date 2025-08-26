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

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.Maplock.Lock()

		for _, user := range this.server.OnlineMap {
			sendmsg := "[" + user.Addr + "]" + user.Name + ":在线...\n"
			this.SendMsg(sendmsg)
		}

		this.server.Maplock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		NewName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[NewName]

		if ok {
			this.SendMsg("当前用户名已被使用")
		} else {
			this.server.Maplock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[NewName] = this
			this.server.Maplock.Unlock()

			this.Name = NewName
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		name := strings.Split(msg, "|")[1]

		if name == "" {
			this.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\" 格式。\n")
			return
		}
		accepter, ok := this.server.OnlineMap[name]
		if !ok {
			this.SendMsg("该用户不存在")
			return
		}
		str := strings.Split(msg, "|")[2]
		if str == "" {
			this.SendMsg("消息为空，请重新输入")
			return
		}
		accepter.SendMsg(this.Name + "对您说： " + str)

	} else {
		this.server.BoardCast(this, msg)
	}

}
