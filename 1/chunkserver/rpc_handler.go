package main



import(
	"fmt"
)

type RpcHandler struct{
}

func (self *RpcHandler) WriteChunk(req *WriteChunkReq, res *WriteChunkRes) error{
	fmt.Printf("write chunk: req[%+v]\n", req);
	return g_chunkserver.WriteChunk(req.Chunk_id, req.Data);
}

func (self *RpcHandler) CreateChunk(req *CreateChunkReq, res *CreateChunkRes) error{
	fmt.Printf("create chunk: req[%+v]\n", req);
	return g_chunkserver.CreateChunk(req.Chunk_id)
}
