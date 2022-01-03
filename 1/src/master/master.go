package main

import (
	"fmt"
	//"time"
	"common"
	"math/rand"
	"sync"
)

type Master struct {
	map_chunk_servers map[uint64]*ChunkserverClient

	map_id_vdi   map[uint32]*Vdi
	map_name_vdi map[string]*Vdi
	vdi_id_seed  uint32
	lock         *sync.Mutex
}

var g_master = Master{
	map_chunk_servers: map[uint64]*ChunkserverClient{},
	map_id_vdi:        map[uint32]*Vdi{},
	map_name_vdi:      map[string]*Vdi{},
	lock:              &sync.Mutex{},
}

func (self *Master) GetVdiById(vdiid uint32) *Vdi {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.map_id_vdi[vdiid]
}
func (self *Master) GetVdiByName(vdiname string) *Vdi {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.map_name_vdi[vdiname]
}
func (self *Master) AddVdi(vdi *Vdi) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.map_id_vdi[vdi.id] = vdi
	self.map_name_vdi[vdi.name] = vdi
}

func (self *Master) select_chunkserver(chunkid uint64) *ChunkserverClient {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.map_chunk_servers) == 0 {
		common.Error("empty chunkserver")
		return nil
	}
	var idx = int(rand.Uint32()) % len(self.map_chunk_servers)
	var i = 0
	for _, chunkserver := range self.map_chunk_servers {
		if i == idx {
			return chunkserver
		}
		i++
	}
	return nil
}
func (self *Master) Heartbeat(chunkserver_id uint64) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	var chunkserver = self.map_chunk_servers[chunkserver_id]
	if chunkserver == nil {
		var ip, port = common.Unpack_chunkserver_id(chunkserver_id)
		self.map_chunk_servers[chunkserver_id] = new_chunkserver_client(chunkserver_id, ip, port)
		common.Info("new chunkserver: [%v:%v]\n", ip, port)
	} else {
	}
	return nil
}

func (self *Master) generate_vdi_id() uint32 {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.vdi_id_seed++
	return self.vdi_id_seed
}

func (self *Master) CreateVdi(req *common.CreateVdiReq, res *common.CreateVdiRes) error {
	var vdi = self.GetVdiByName(req.Vdi_name)
	if vdi != nil {
		return fmt.Errorf("conflict vdi name")
	}

	vdi = new_vdi(req.Vdi_name, self.generate_vdi_id(), req.Vdi_size)
	var err = vdi.create_chunk()
	if err != nil {
		common.Error("failed to create chunk: vdiname[%v] vdiid[%v]", vdi.name, vdi.id)
		return err
	}
	self.AddVdi(vdi)

	res.Vdi_name = vdi.name
	res.Vdi_id = vdi.id
	return nil
}

func (self *Master) GetChunkMeta(chunkid uint64) uint64 {
	var vdiid = common.Get_vdiid_from_chunkid(chunkid)
	var vdi = self.map_id_vdi[vdiid]
	if vdi == nil {
		common.Error("unknown vdiid[%v]\n", vdiid)
		return 0
	} else {
		return vdi.get_chunkserverid_by_chunkid(chunkid)
	}
}
func (self *Master) GetVdiMeta(vdiname string, vdi_id *uint32, vdi_size *uint64) error {
	var vdi = self.GetVdiByName(vdiname)
	if vdi == nil {
		common.Error("unknown vdiname[%v]\n", vdiname)
		return fmt.Errorf("unknown vdiname[%v]\n", vdiname)
	}
	*vdi_id = vdi.id
	*vdi_size = vdi.size
	return nil
}
