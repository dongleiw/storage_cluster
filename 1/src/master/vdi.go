package main

import (
	"common"
	"fmt"
	//"time"
	"sync"
)

type Chunk struct {
	chunk_id uint64
}
type Vdi struct {
	name                      string
	id                        uint32
	size                      uint64
	map_chunkid_chunkserverid map[uint64]uint64
	lock                      *sync.Mutex
}

func new_vdi(vdi_name string, vdi_id uint32, size uint64) *Vdi {
	return &Vdi{
		name:                      vdi_name,
		id:                        vdi_id,
		size:                      size,
		map_chunkid_chunkserverid: map[uint64]uint64{},
	}
}

func (self *Vdi) create_chunk() error {
	// 创建chunk
	var i uint64 = 0
	for i = 0; i < self.size/256; i++ {
		var chunkid = common.Generate_chunkid(self.id, uint32(i))
		var chunk_server = g_master.select_chunkserver(chunkid)
		if chunk_server == nil {
			return fmt.Errorf("failed to select chunkserver")
		}

		if err := chunk_server.create_chunk(chunkid); err == nil {
			self.map_chunkid_chunkserverid[chunkid] = chunk_server.id
		} else {
			common.Error("failed to create chunk[%v]: %v\n", chunkid, err)
			return fmt.Errorf("failed to create chunk[%v]", chunkid)
		}
	}
	return nil
}

func (self *Vdi) get_chunkserverid_by_chunkid(chunkid uint64) uint64 {
	return self.map_chunkid_chunkserverid[chunkid]
}
