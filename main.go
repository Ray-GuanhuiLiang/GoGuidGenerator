package main

import (
	"fmt"
	"github.com/Ray-GuanhuiLiang/GoGuidGenerator/guid"
	"time"
)

func main() {
	test1()
}

func test1() {
	g, err := guid.NewGuid(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:=0; i<3; i++ {
		fmt.Println(g.Generate())
	}
}

func test2() {
	g, err := guid.NewGuid(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:=0; i<3; i++ {
		g.Generate()
	}
	t1:= time.Now()
	for i:=0; i<1000000; i++ {
		g.Generate()
	}
	t2:= time.Now()
	fmt.Println(t2.Sub(t1))
}