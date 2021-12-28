package main

type Chunk struct{
	chunk_id int64
}
type Vdi struct{
	name string
	id int32
	size int64
	map_vdi_chunks map[int64]*Chunk
}

func new_vdi(vdi_name string, vdi_id int32, size int64)*Vdi{
	return &Vdi{
		name:vdi_name,
		id: vdi_id,
		size: size,
		map_vdi_chunks: map[int64]*Chunk{},
	}
}
