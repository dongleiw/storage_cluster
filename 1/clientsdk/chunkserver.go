package main

import(
	"io/ioutil"
	"fmt"
)
type Chunkserver struct{
	datadir string
}
var g_chunkserver = Chunkserver{
	datadir:"./data",
}

func (self *Chunkserver) WriteChunk(chunkid int64, data []byte) error{
	return ioutil.WriteFile( fmt.Sprintf("%v/%v", self.datadir, chunkid), data, 600);
}
