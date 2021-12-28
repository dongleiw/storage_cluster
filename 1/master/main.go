package main

import(
	"net"
	"fmt"
	"net/rpc"
)
func main(){
	fmt.Println("server start")
	rpc.Register(&RpcHandler{})
	var l, err = net.Listen("tcp", ":2345")
	if err != nil {
	    panic(err)
	}
	rpc.Accept(l)
}
