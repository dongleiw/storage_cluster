package main

type Chunk struct{
	chunk_id uint64
}
type Vdi struct{
	name string
	id uint32
	size uint64
	map_vdi_chunks map[uint64]*Chunk
}

func new_vdi(vdi_name string, vdi_id uint32, size uint64)*Vdi{
	return &Vdi{
		name:vdi_name,
		id: vdi_id,
		size: size,
		map_vdi_chunks: map[uint64]*Chunk{},
	}
}
