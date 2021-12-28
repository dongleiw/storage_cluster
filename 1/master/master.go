package main

import(
	"fmt"
)

type Chunkserver struct{
	id uint64
	ip string
	port int
	rpcclient *rpc.Client
}

type Master struct{
	map_chunk_servers map[uint64]*Chunkserver

	map_id_vdi map[uint32]*Vdi
	map_name_vdi map[string]*Vdi
	vdi_id_seed uint32
}

var g_master = Master{
	map_id_vdi: map[uint32]*Vdi{},
	map_name_vdi: map[string]*Vdi{},
};

// 选择一个chunkserver创建chunk
func (self *Master) create_chunk(chunkid uint64) error{
}
func (self *Master) Heartbeat(req *HeartbeatReq) error{
	var chunkserver = self.map_chunk_servers[req.Chunk_server_id]
	if chunkserver==nil{
		chunkserver = &Chunkserver{
			id: req.Chunk_server_id,
			ip: req.Ip,
			port: req.Port,
		}
		self.map_chunk_servers[req.Chunk_server_id] = chunkserver
	}
	return nil
}

// TODO lock
func (self *Master) CreateVdi(req *CreateVdiReq, res *CreateVdiRes) error{
	var _, exist = self.map_name_vdi[req.Vdi_name]
	if exist{
		return fmt.Errorf("conflict vdi name");
	}

	self.vdi_id_seed++
	var vdi = new_vdi(req.Vdi_name, self.vdi_id_seed, req.Vdi_size)

	// 创建chunk
	var i uint64 = 0;
	for i=0; i<req.Vdi_size/256; i++ {
		var chunkid uint64 = uint64(vdi.id) <<32 | i;
		self.create_chunk(chunkid)
	}


	self.map_id_vdi[vdi.id] = vdi
	self.map_name_vdi[vdi.name] = vdi

	res.Vdi_name = vdi.name
	res.Vdi_id = vdi.id
	return nil
}
