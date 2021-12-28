package main

import(
	"fmt"
)

type Master struct{
	map_id_vdi map[int32]*Vdi
	map_name_vdi map[string]*Vdi
	vdi_id_seed int32
}

var g_master = Master{
	map_id_vdi: map[int32]*Vdi{},
	map_name_vdi: map[string]*Vdi{},
};

func (self *Master) CreateVdi(req *CreateVdiReq, res *CreateVdiRes) error{
	var _, exist = self.map_name_vdi[req.Vdi_name]
	if exist{
		return fmt.Errorf("conflict vdi name");
	}

	self.vdi_id_seed++
	var vdi = new_vdi(req.Vdi_name, self.vdi_id_seed, req.Vdi_size)

	self.map_id_vdi[vdi.id] = vdi
	self.map_name_vdi[vdi.name] = vdi

	res.Vdi_name = vdi.name
	res.Vdi_id = vdi.id
	return nil
}
