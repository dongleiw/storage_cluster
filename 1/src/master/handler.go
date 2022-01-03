package main

import (
	//"fmt"
	"common"
)

type RpcHandler struct {
}

func (self *RpcHandler) Heartbeat(req *common.HeartbeatReq, res *common.HeartbeatRes) error {
	//fmt.Printf("heartbeat: req[%+v]\n", req);
	return g_master.Heartbeat(req.Chunk_server_id)
}
func (self *RpcHandler) CreateVdi(req *common.CreateVdiReq, res *common.CreateVdiRes) error {
	common.Debug("create vdi: req[%+v]\n", req)
	return g_master.CreateVdi(req, res)
}
func (self *RpcHandler) GetChunkMeta(req *common.GetChunkMetaReq, res *common.GetChunkMetaRes) error {
	common.Debug("get chunk meta: req[%+v]\n", req)
	var chunk_server_id = g_master.GetChunkMeta(req.Chunk_id)

	res.Chunkserver_id = chunk_server_id
	return nil
}
func (self *RpcHandler) GetVdiMeta(req *common.GetVdiMetaReq, res *common.GetVdiMetaRes) error {
	common.Debug("get vdi chunk meta: req[%+v]\n", req)
	var err = g_master.GetVdiMeta(req.Vdi_name, &res.Vdi_id, &res.Vdi_size)

	return err
}
