package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type User struct {
	ServerIp   string
	ServerPort int
	Conn       net.Conn
}

var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "set default ip")
	flag.IntVar(&port, "port", 8090, "set default port")
}

func NewUser(ip string, port int) *User {
	user := &User{
		ServerIp:   ip,
		ServerPort: port,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil
	} else {
		user.Conn = conn
		return user
	}
}

func (user *User) Menu() {
	fmt.Println("1. enter 1 to boardcast")
	fmt.Println("2. enter 2 to private chat")
	fmt.Println("3. enter 3 to change name")
	fmt.Println("4. enter 0 to exit")
}

func (user *User) UpdateName() {
	fmt.Println("enter the name to replace original name")
	var name string
	fmt.Scan(&name)
	msg := "rename|" + name
	_, err := user.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("fail to rename!")
	}
}

func (user *User) Boardcast() {
	fmt.Println("enter the message you want to boardcast:")
	var msg string
	fmt.Scan(&msg)
	_, err := user.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("fail to boardcast!")
	}
}

func (user *User) ReceiveMsg() {
	io.Copy(os.Stdout, user.Conn)
}

func main() {
	flag.Parse()
	user := NewUser(ip, port)
	if user == nil {
		fmt.Println("fail to connect!")
	} else {
		fmt.Println("connect successfully")
	}
	go user.ReceiveMsg()
	var str string = "!!!"
	for {
		user.Menu()
		fmt.Scan(&str)
		fmt.Printf("str === %ss", str)

		switch str {
		case "0":
			fmt.Println("exit")
			return
		case "1":
			user.Boardcast()
		case "2":
			fmt.Println("private chat")
		case "3":
			user.UpdateName()
		default:
			fmt.Println("wrong input!")
		}
	}
}
