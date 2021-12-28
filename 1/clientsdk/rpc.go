package main

type WriteChunkReq struct{
	Chunk_id int64
	Data []byte
}
type WriteChunkRes struct{
	Chunk_id int64
}
