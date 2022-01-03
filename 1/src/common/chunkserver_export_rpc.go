package common

type WriteChunkReq struct {
	Chunk_id uint64
	Data     []byte
}
type WriteChunkRes struct {
	Chunk_id uint64
}

type ReadChunkReq struct {
	Chunk_id uint64
}
type ReadChunkRes struct {
	Data []byte
}

type CreateChunkReq struct {
	Chunk_id uint64
}
type CreateChunkRes struct {
	Chunk_id uint64
}
