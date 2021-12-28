package main


type HeartbeatReq struct{
	Chunk_server_id uint64
	Ip string
	Port int
}
type HeartbeatRes struct{
}

type CreateVdiReq struct{
	Vdi_name string
	Vdi_size uint64
}
type CreateVdiRes struct{
	Vdi_name string
	Vdi_id uint32
}
