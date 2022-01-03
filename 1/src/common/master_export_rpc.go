package common

type HeartbeatReq struct {
	Chunk_server_id uint64
}
type HeartbeatRes struct {
}

type CreateVdiReq struct {
	Vdi_name string
	Vdi_size uint64
}
type CreateVdiRes struct {
	Vdi_name string
	Vdi_id   uint32
}

type GetChunkMetaReq struct {
	Chunk_id uint64
}
type GetChunkMetaRes struct {
	Chunkserver_id uint64
}

type GetVdiMetaReq struct {
	Vdi_name string
}
type GetVdiMetaRes struct {
	Vdi_id   uint32
	Vdi_size uint64
}
