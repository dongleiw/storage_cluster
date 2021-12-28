package main

import(
	"fmt"
)
type RpcHandler struct{
}

func (self *RpcHandler) Heartbeat(req *HeartbeatReq, res *HeartbeatRes) error{
	fmt.Printf("heartbeat: req[%+v]\n", req);
	return g_master.HeartBeat(req)
}
func (self *RpcHandler) CreateVdi(req *CreateVdiReq, res *CreateVdiRes) error{
	fmt.Printf("create vdi: req[%+v]\n", req);
	return g_master.CreateVdi(req, res);
}
