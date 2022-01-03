package main

import (
	"common"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"time"
)

type Chunkserver struct {
	chunk_server_id   uint64
	map_chunkid_chunk map[uint64]bool
	datadir           string
	ip                string
	port              int
	master_ip         string
	master_port       int
}

var g_chunkserver *Chunkserver = nil

func init_chunkserver(ip string, port int, master_ip string, master_port int) {
	g_chunkserver = &Chunkserver{
		datadir:           "./data",
		chunk_server_id:   common.Generate_chunkserver_id(ip, port),
		map_chunkid_chunk: map[uint64]bool{},
		ip:                ip,
		port:              port,
		master_ip:         master_ip,
		master_port:       master_port,
	}
}

func (self *Chunkserver) WriteChunk(chunkid uint64, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%v/%v", self.datadir, chunkid), data, 0600)
}
func (self *Chunkserver) ReadChunk(chunkid uint64) ([]byte, error) {
	var content, err = ioutil.ReadFile(fmt.Sprintf("%v/%v", self.datadir, chunkid))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (self *Chunkserver) CreateChunk(chunkid uint64) error {
	if _, e := self.map_chunkid_chunk[chunkid]; e {
		return fmt.Errorf("chunkid[%v] already exists", chunkid)
	} else {
		self.map_chunkid_chunk[chunkid] = true
		return nil
	}
}

func (self *Chunkserver) Heartbeat() {
	var client, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", self.master_ip, self.master_port))
	if err != nil {
		common.Error("failed to connect to master: %v\n", err)
	}

	for {
		if client == nil {
			client, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", self.master_ip, self.master_port))
			if err != nil {
				common.Error("failed to connect to master: %v\n", err)
				client = nil
			}
		}
		time.Sleep(time.Second * 1)

		if client != nil {
			var req = &common.HeartbeatReq{
				Chunk_server_id: self.chunk_server_id,
			}
			var res = &common.HeartbeatRes{}
			err = client.Call("RpcHandler.Heartbeat", req, res)
			if err != nil {
				common.Error("failed to heartbeat: %v", err)
				client = nil
			} else {
				common.Debug("heartbeat: %v %v", req, res)
			}
		}
	}
}
