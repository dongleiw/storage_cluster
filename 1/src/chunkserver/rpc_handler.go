package main

import (
	//"fmt"
	"common"
)

type RpcHandler struct {
}

func (self *RpcHandler) WriteChunk(req *common.WriteChunkReq, res *common.WriteChunkRes) error {
	//common.Debug("write chunk: chunkid[%v]\n", req.Chunk_id)
	return g_chunkserver.WriteChunk(req.Chunk_id, req.Data)
}
func (self *RpcHandler) ReadChunk(req *common.ReadChunkReq, res *common.ReadChunkRes) error {
	//common.Debug("read chunk: chunkid[%v]\n", req.Chunk_id)
	var content, err = g_chunkserver.ReadChunk(req.Chunk_id)
	if err != nil {
		return err
	}
	res.Data = content
	return nil
}

func (self *RpcHandler) CreateChunk(req *common.CreateChunkReq, res *common.CreateChunkRes) error {
	common.Debug("create chunk: req[%+v]\n", req)
	return g_chunkserver.CreateChunk(req.Chunk_id)
}
