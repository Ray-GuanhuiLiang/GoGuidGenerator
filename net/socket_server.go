package server

import (
	"github.com/Ray-GuanhuiLiang/GoGuidGenerator/common"
	"log"
	"net"
	"strconv"
)

type TcpServer struct {
	generator common.Generator
	exit      chan int
}

func NewTcpServer(generator common.Generator) *TcpServer {
	e := make(chan int)
	return &TcpServer{generator, e}
}

func (this *TcpServer) Start() error {
	listener, err := net.Listen("tcp", ":5588")
	if err != nil {
		log.Fatal(err)
		return err
	}
	go this.handleAccept(listener)
	return nil
}

func (this *TcpServer) Wait() {
	<-this.exit
}

func (this *TcpServer) handleAccept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClient(conn)
	}
}

func (this *TcpServer) handleClient(conn net.Conn) {
	buf := make([]byte, 1024)
	log.Printf("Accept %s", conn)
	for {
		n, err := conn.Read(buf)
		log.Println(n)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		if n <= 0 {
			conn.Close()
			return
		}
		guid, err := this.generator.Generate()
		if err != nil {
			_, err = conn.Write([]byte("0\n"))
		} else {
			_, err = conn.Write([]byte(strconv.FormatUint(guid, 10) + "\n"))
		}
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
	}

}
