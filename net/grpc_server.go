package server

import (
	"github.com/Ray-GuanhuiLiang/GoGuidGenerator/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	generator common.Generator
	serv      *grpc.Server
	exit      chan int
}

func NewGrpcServer(generator common.Generator) *GrpcServer {
	e := make(chan int)
	serv := grpc.NewServer()
	s := &GrpcServer{generator, serv, e}
	RegisterGuidServer(serv, s)
	return s
}

func (this *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", ":5588")
	if err != nil {
		log.Fatal(err)
		return err
	}
	go this.serv.Serve(listener)
	return nil
}

func (this *GrpcServer) Wait() {
	<-this.exit
}

func (this *GrpcServer) GetGuid(context.Context, *Req) (*Resp, error) {
	id, err := this.generator.Generate()
	var (
		r    Resp
		code int32
		guid uint64
	)
	if err != nil {
		code = 1
		guid = 0
	} else {
		code = 0
		guid = id
	}
	r.Code = &code
	r.Guid = &guid
	return &r, err
}
