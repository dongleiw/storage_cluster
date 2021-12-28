package main

import(
	"io/ioutil"
	"fmt"
)
type Chunkserver struct{
	map_chunkid_chunk map[int64]bool
	datadir string
}
var g_chunkserver = Chunkserver{
	datadir:"./data",
	map_chunkid_chunk: map[int64]bool{},
}

func (self *Chunkserver) WriteChunk(chunkid int64, data []byte) error{
	return ioutil.WriteFile( fmt.Sprintf("%v/%v", self.datadir, chunkid), data, 600);
}

func (self *Chunkserver) CreateChunk(chunkid int64) error{
	if _,e:=self.map_chunkid_chunk[chunkid]; e{
		return fmt.Errorf("chunkid[%v] already exists", chunkid);
	}else{
		self.map_chunkid_chunk[chunkid] = true
		return nil
	}
}
