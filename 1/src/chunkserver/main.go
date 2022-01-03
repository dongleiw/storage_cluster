package main

import (
	"common"
	"fmt"
	"net"
	"net/rpc"
)

const IP = "127.0.0.1"
const PORT = 2346
const MASTER_IP = "127.0.0.1"
const MASTER_PORT = 2345

func listen() {
	common.Info("server start")
	rpc.Register(&RpcHandler{})
	var l, err = net.Listen("tcp", fmt.Sprintf(":%v", PORT))
	if err != nil {
		panic(err)
	}
	rpc.Accept(l)
}
func main() {
	common.Log_init()

	init_chunkserver(IP, PORT, MASTER_IP, MASTER_PORT)
	go g_chunkserver.Heartbeat()

	listen()
}
