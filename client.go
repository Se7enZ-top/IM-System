package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("连接错误：", err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) Menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字<<<<<")
		return false
	}

}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "默认ip为127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "默认端口为8888")

}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) PublicTalk() {
	fmt.Println("欢迎进入公聊模式")
	var chatStr string

	fmt.Println(">>>请输入要发送的内容（exit退出）<<<")
	fmt.Scanln(&chatStr)

	for chatStr != "exit" {
		if len(chatStr) != 0 {
			sendMsg := chatStr + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("err :", err)
				break
			}
		}

		chatStr = ""
		fmt.Scanln(&chatStr)

	}

}

func (client *Client) FindOnlineUser() {
	sendMsg := "who\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
}

func (client *Client) PrivateTalk() {
	fmt.Println("欢迎进入私聊模式")

	var remoteName string
	var chatMsg string

	client.FindOnlineUser()

	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容, exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}

		client.FindOnlineUser()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}

}
func (client *Client) Rename() bool {
	fmt.Println(">>>>>请输入用户名：<<<<<")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("write err:", err)
		return false
	}
	return true

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		switch client.flag {
		case 1:
			client.PublicTalk()
			break
		case 2:
			client.PrivateTalk()
			break
		case 3:
			client.Rename()
			break
		}
	}
	return
}

func main() {
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>连接失败<<<<<")
		return
	}
	fmt.Println(">>>>>连接成功<<<<<")

	go client.DealResponse()

	client.Run()
}
