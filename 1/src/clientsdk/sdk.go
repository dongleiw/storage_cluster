package clientsdk

import (
	"fmt"
	"net/rpc"
)

type Vdi struct {
	Vdi_name string
	Vdi_id   uint32
	Vdi_size uint64

	offset         uint64
	chunk_location []uint64 // chunkidx->chunkserverid
}

var (
	SEEK_SET = 1
	SEEK_CUR = 2
)

func Init(master_ip string, master_port int) {
	g_clientsdk = &ClientSDK{
		master_ip:                   master_ip,
		master_port:                 master_port,
		map_chunkserverid_rpcclient: map[uint64]*rpc.Client{},
	}
}

func Open(vdiname string) (*Vdi, error) {
	return g_clientsdk.open_vdi(vdiname)
}

func Read(vdi *Vdi, len uint64) ([]byte, error) {
	return g_clientsdk.read_vdi(vdi, len)
}

func Write(vdi *Vdi, data []byte) (uint64, error) {
	return g_clientsdk.write_vdi(vdi, data)
}

func Seek(vdi *Vdi, offset uint64, whence int) error {
	if whence == SEEK_SET {
		if offset >= vdi.Vdi_size {
			return fmt.Errorf("invalid offset")
		}
		vdi.offset = offset
	} else if whence == SEEK_CUR {
		if vdi.offset+offset >= vdi.Vdi_size {
			return fmt.Errorf("invalid offset")
		}
		vdi.offset += offset
	}
	return fmt.Errorf("unknown seek type")
}

func Close(vdi *Vdi) error {
	vdi.chunk_location = nil
	return nil
}
