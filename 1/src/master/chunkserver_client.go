package main

import (
	"common"
	"fmt"
	"net/rpc"
	"sync"
)

type ChunkserverClient struct {
	id        uint64
	ip        string
	port      int
	rpcclient *rpc.Client
	lock      *sync.Mutex
}

func new_chunkserver_client(id uint64, ip string, port int) *ChunkserverClient {
	return &ChunkserverClient{
		id:        id,
		ip:        ip,
		port:      port,
		rpcclient: nil,
		lock:      &sync.Mutex{},
	}
}

// rpc保持长连接?
// rpc.call thready safety
// rpc有自动重连么?
func (self *ChunkserverClient) create_chunk(chunkid uint64) error {
	self.lock.Lock()
	if self.rpcclient == nil {
		var client, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", self.ip, self.port))
		if err != nil {
			self.lock.Unlock()
			return fmt.Errorf("failed to connect to chunkserver[%v:%v]: %v\n", self.ip, self.port, err)
		}
		self.rpcclient = client
	}
	self.lock.Unlock()
	var req = &common.CreateChunkReq{
		Chunk_id: chunkid,
	}
	var res = &common.CreateChunkRes{}
	return self.rpcclient.Call("RpcHandler.CreateChunk", req, res)
}
