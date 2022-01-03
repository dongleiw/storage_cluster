package main

import (
	//"fmt"
	"common"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

func main() {
	common.Log_init()
	common.Info("server start")

	rand.Seed(time.Now().Unix())

	rpc.Register(&RpcHandler{})
	var l, err = net.Listen("tcp", ":2345")
	if err != nil {
		panic(err)
	}
	rpc.Accept(l)
}
