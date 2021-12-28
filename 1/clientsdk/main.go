package main

import(
	"fmt"
	"net/rpc"
	"net"
	"time"
)

var DATA_PORT = 2346
func listen(){
	fmt.Println("server start")
	rpc.Register(&RpcHandler{})
	var l, err = net.Listen("tcp", fmt.Sprintf(":%v", DATA_PORT))
	if err != nil {
	    panic(err)
	}
	rpc.Accept(l)
}
func heartbeat(){
	var client, err = rpc.Dial("tcp", ":2345")
	if err != nil {
		fmt.Printf("failed to connect to master: %v\n", err)
	}

	for{
		if(client==nil){
			client, err = rpc.Dial("tcp", ":2345")
			if err != nil {
				fmt.Printf("failed to connect to master: %v\n", err)
			}
		}
		time.Sleep(time.Second*1);

		if(client!=nil){
			var res = &HeartbeatRes{}
			err = client.Call("RpcHandler.Heartbeat", HeartbeatReq{A: 1}, res)
			fmt.Println(res, err)
		}
	}
}
func main(){
	go heartbeat()

	listen()
}
