package main

import(
	"fmt"
	"net/rpc"
	"net"
	"time"
)

const DATA_PORT = 2346
const CHUNK_SERVER_ID = 1
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
			var req = &HeartbeatReq{
				Chunk_server_id: CHUNK_SERVER_ID,
			}
			var res = &HeartbeatRes{}
			err = client.Call("RpcHandler.Heartbeat", req, res)
			fmt.Println(res, err)
		}
	}
}
func main(){
	go heartbeat()

	listen()
}
