package clientsdk

import (
	"bytes"
	"common"
	"fmt"
	"net/rpc"
)

const CHUNK_SIZE = 256

type ClientSDK struct {
	master_ip                   string
	master_port                 int
	master_rpcclient            *rpc.Client
	map_chunkserverid_rpcclient map[uint64]*rpc.Client
}

var (
	g_clientsdk *ClientSDK = nil
)

func (self *ClientSDK) get_meta_of_vdi(vdiname string, vdi *Vdi) error {
	if self.master_rpcclient == nil {
		var master_rpcclient, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", self.master_ip, self.master_port))
		if err != nil {
			return fmt.Errorf("failed to connect to master: %v", err)
		}
		self.master_rpcclient = master_rpcclient
	}

	var req = &common.GetVdiMetaReq{
		Vdi_name: vdiname,
	}
	var res = &common.GetVdiMetaRes{}
	var err = self.master_rpcclient.Call("RpcHandler.GetVdiMeta", req, res)
	if err != nil {
		self.master_rpcclient = nil
		return fmt.Errorf("failed to GetVdiMeta: %v", err)
	} else {
		common.Debug("success to GetVdiMeta: %v %v", req, res)
		vdi.Vdi_id = res.Vdi_id
		vdi.Vdi_size = res.Vdi_size
		return nil
	}
}
func (self *ClientSDK) call_master_rpc(method string, req interface{}, res interface{}) error {
	var err = self.master_rpcclient.Call(method, req, res)
	if err != nil {
		self.master_rpcclient = nil
	}
	return err
}
func (self *ClientSDK) open_vdi(vdiname string) (*Vdi, error) {
	var vdi = &Vdi{
		Vdi_name: vdiname,
	}
	var err = g_clientsdk.get_meta_of_vdi(vdiname, vdi)
	if err != nil {
		common.Error("open failed: vdiname[%v] vdiid[%v] %v", vdi.Vdi_name, vdi.Vdi_id, err)
		return nil, err
	}
	vdi.chunk_location = make([]uint64, vdi.Vdi_size/CHUNK_SIZE, vdi.Vdi_size/CHUNK_SIZE)
	common.Debug("open succ: vdiname[%v] vdiid[%v]", vdi.Vdi_name, vdi.Vdi_id)
	return vdi, err
}
func (self *ClientSDK) write_chunk(vdi *Vdi, chunkid uint64, data []byte) error {
	var err error
	// get meta of chunkid
	var chunkidx = uint32(chunkid & 0xffffffff)
	var chunkserver_id = vdi.chunk_location[chunkidx]
	if chunkserver_id == 0 {
		var req = &common.GetChunkMetaReq{
			Chunk_id: chunkid,
		}
		var res = &common.GetChunkMetaRes{}
		err = self.call_master_rpc("RpcHandler.GetChunkMeta", req, res)
		if err != nil {
			common.Error("failed to GetChunkMeta: %v", err)
			return err
		} else {
			//Info("success to GetChunkMeta: %v %v", req, res)
			chunkserver_id = res.Chunkserver_id
			vdi.chunk_location[chunkidx] = chunkserver_id
		}
	}

	// get rpc connection
	var chunkserver_rpcclient = self.map_chunkserverid_rpcclient[chunkserver_id]
	if chunkserver_rpcclient == nil {
		var chunk_server_ip, chunk_server_port = common.Unpack_chunkserver_id(chunkserver_id)
		chunkserver_rpcclient, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", chunk_server_ip, chunk_server_port))
		if err != nil {
			common.Error("failed to connect to master: %v\n", err)
			return err
		}
		self.map_chunkserverid_rpcclient[chunkserver_id] = chunkserver_rpcclient
	}

	var req = &common.WriteChunkReq{
		Chunk_id: chunkid,
		Data:     data,
	}
	var res = &common.WriteChunkRes{}
	err = chunkserver_rpcclient.Call("RpcHandler.WriteChunk", req, res)
	if err != nil {
		delete(self.map_chunkserverid_rpcclient, chunkserver_id)
		common.Error("failed to WriteChunk: %v", err)
		return err
	}
	return nil
}
func (self *ClientSDK) read_chunk(vdi *Vdi, chunkid uint64) ([]byte, error) {
	var err error
	// get meta of chunkid
	var chunkidx = uint32(chunkid & 0xffffffff)
	var chunkserver_id = vdi.chunk_location[chunkidx]
	if chunkserver_id == 0 {
		var req = &common.GetChunkMetaReq{
			Chunk_id: chunkid,
		}
		var res = &common.GetChunkMetaRes{}
		err = self.call_master_rpc("RpcHandler.GetChunkMeta", req, res)
		if err != nil {
			common.Error("failed to GetChunkMeta: %v", err)
			return nil, err
		} else {
			//Info("success to GetChunkMeta: %v %v", req, res)
			chunkserver_id = res.Chunkserver_id
			vdi.chunk_location[chunkidx] = chunkserver_id
		}
	}

	// get rpc connection
	var chunkserver_rpcclient = self.map_chunkserverid_rpcclient[chunkserver_id]
	if chunkserver_rpcclient == nil {
		var chunk_server_ip, chunk_server_port = common.Unpack_chunkserver_id(chunkserver_id)
		chunkserver_rpcclient, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", chunk_server_ip, chunk_server_port))
		if err != nil {
			common.Error("failed to connect to master: %v\n", err)
			return nil, err
		}
		self.map_chunkserverid_rpcclient[chunkserver_id] = chunkserver_rpcclient
	}

	var req = &common.ReadChunkReq{
		Chunk_id: chunkid,
	}
	var res = &common.ReadChunkRes{}
	err = chunkserver_rpcclient.Call("RpcHandler.ReadChunk", req, res)
	if err != nil {
		delete(self.map_chunkserverid_rpcclient, chunkserver_id)
		common.Error("failed to WriteChunk: %v", err)
		return nil, err
	}
	return res.Data, nil
}

// 返回成功写入的数据长度, 错误信息
func (self *ClientSDK) write_vdi(vdi *Vdi, data []byte) (uint64, error) {
	if len(data)%CHUNK_SIZE != 0 || vdi.offset+uint64(len(data)) >= vdi.Vdi_size {
		return 0, fmt.Errorf("invalid data len")
	}
	// 串行算chunkid, 获取chunkid所在的chunkserver, 向chunkserver发写数据请求
	var old_offset = vdi.offset
	var begidx = old_offset / CHUNK_SIZE
	var endidx = (old_offset + uint64(len(data))) / CHUNK_SIZE
	for i := begidx; i < endidx; i++ {
		var chunkid = common.Generate_chunkid(vdi.Vdi_id, uint32(i))
		var err = self.write_chunk(vdi, chunkid, data[(i-begidx)*CHUNK_SIZE:(i+1-begidx)*CHUNK_SIZE])
		if err != nil {
			return vdi.offset - old_offset, err
		}
		vdi.offset += CHUNK_SIZE
	}
	return vdi.offset - old_offset, nil
}

// 读取指定长度的数据, 如果提前遇到EOF. 则返回读取的数据, error=nil
func (self *ClientSDK) read_vdi(vdi *Vdi, readlen uint64) ([]byte, error) {
	if readlen%CHUNK_SIZE != 0 {
		return nil, fmt.Errorf("invalid data len")
	}
	var buf = bytes.NewBuffer(nil)
	var old_offset = vdi.offset
	var begidx = old_offset / CHUNK_SIZE
	var endidx uint64 = 0
	if old_offset+uint64(readlen) >= vdi.Vdi_size {
		endidx = vdi.Vdi_size / CHUNK_SIZE
	} else {
		endidx = (old_offset + uint64(readlen)) / CHUNK_SIZE
	}
	for i := begidx; i < endidx; i++ {
		var chunkid = common.Generate_chunkid(vdi.Vdi_id, uint32(i))
		var chunk_data, err = self.read_chunk(vdi, chunkid)
		if err != nil {
			if buf.Len() > 0 {
				return buf.Bytes(), nil
			} else {
				return nil, err
			}
		}
		buf.Write(chunk_data)
		vdi.offset += CHUNK_SIZE
	}
	return buf.Bytes(), nil
}
