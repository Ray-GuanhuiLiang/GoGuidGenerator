package main

import (
	"github.com/Ray-GuanhuiLiang/GoGuidGenerator/guid"
	"github.com/Ray-GuanhuiLiang/GoGuidGenerator/net"
	"fmt"
)

func main() {
	g, err := guid.NewGuid()
	if err != nil {
		fmt.Println(err)
		return
	}
	s := server.NewGrpcServer(g)
	s.Start()
	s.Wait()
	fmt.Println("server end")
}
