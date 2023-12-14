package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	OnlineMap map[string]*Client
	Message   chan string
	mapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server {
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*Client),
		Message: make(chan string),
	}
	return server
}

func (server *Server) ListenChan() {
	for {
		msg := <- server.Message
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.Channel <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) WriteMegToChan(client *Client, msg string) {
	message := "[" + client.Addr + "] " + client.Name + "; " + msg
	server.Message <- message
}

func (server *Server) CheckMsgType(msg string, conn net.Conn, client *Client) {
	if msg == "who" {
		server.mapLock.Lock()
		for _, client := range server.OnlineMap {
			message := "[" + client.Addr + "] " + client.Name + ";\n"
			conn.Write([]byte(message))
		}
		server.mapLock.Unlock()
		
	} else if strings.HasPrefix(msg, "rename") {
		c := server.OnlineMap[client.Name]
		delete(server.OnlineMap, client.Name)
		client.Name = msg[7:]
		server.OnlineMap[client.Name] = c
		message := fmt.Sprintf("changed name to %s", client.Name)
		server.WriteMegToChan(client, message)
	} else {
		server.WriteMegToChan(client, msg)
	}
}

func (server *Server) ReadMsgFromClient(client *Client, conn net.Conn, isLive chan bool) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			delete(server.OnlineMap, client.Name)
			msg := "offline!"
			server.WriteMegToChan(client, msg)
			client = nil
			return
		} else if err != nil && err != io.EOF {
			fmt.Println("read error!")
			return
		} else {
			msg := string(buf[:n-1])
			server.CheckMsgType(msg, conn, client)
			isLive <- true
		}
	}
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("success to connect!")
	
	client := NewClient(conn)
	isLive := make(chan bool)
	go server.ReadMsgFromClient(client, conn, isLive)

	server.mapLock.Lock()
	server.OnlineMap[client.Name] = client
	server.mapLock.Unlock()

	server.WriteMegToChan(client, "online!")

	for {
		select {
		case <- isLive:

		case <- time.After(time.Minute * 60):
			conn.Write([]byte("you're kick out"))
			close(client.Channel)
			conn.Close()
			return
		}
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("listen error!")
	}

	go server.ListenChan()

	defer listener.Close()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error!")
			continue
		}
		go server.Handler(conn)
	}
}
