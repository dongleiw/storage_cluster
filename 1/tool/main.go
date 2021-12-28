package main

import(
	"fmt"
	"net/rpc"
)
const MASTER_PORT = 2345
func main(){
	var client, err = rpc.Dial("tcp", fmt.Sprintf(":%v", MASTER_PORT))
	if err != nil {
		fmt.Printf("failed to connect to master: %v\n", err)
	}

	if(client!=nil){
		var req = &CreateVdiReq{
			Vdi_name: "hello",
			Vdi_size: 1024*1024*10,
		}
		var res = &CreateVdiRes{}
		err = client.Call("RpcHandler.CreateVdi", req, res)
		fmt.Println(res, err)
	}
}
