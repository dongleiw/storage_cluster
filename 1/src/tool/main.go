package main

import (
	"clientsdk"
	"common"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

const MASTER_IP = "127.0.0.1"
const MASTER_PORT = 2345

func create_vdi(vdiname string, vdisize uint64) {

	var client, err = rpc.Dial("tcp", fmt.Sprintf(":%v", MASTER_PORT))
	if err != nil {
		common.Error("failed to connect to master: %v\n", err)
		return
	}

	if client != nil {
		var req = &common.CreateVdiReq{
			Vdi_name: vdiname,
			Vdi_size: vdisize,
		}
		var res = &common.CreateVdiRes{}
		err = client.Call("RpcHandler.CreateVdi", req, res)
		fmt.Println(res, err)
	}
}

func get_chunkid_meta(args []string) {
	var chunkid, _ = strconv.ParseUint(args[0], 10, 64)

	var client, err = rpc.Dial("tcp", fmt.Sprintf(":%v", MASTER_PORT))
	if err != nil {
		common.Error("failed to connect to master: %v\n", err)
		return
	}

	var req = &common.GetChunkMetaReq{
		Chunk_id: chunkid,
	}
	var res = &common.GetChunkMetaRes{}
	err = client.Call("RpcHandler.GetChunkMeta", req, res)
	if err != nil {
		common.Error("failed to GetChunkMeta: %v", err)
	} else {
		common.Info("success to GetChunkMeta: %v %v", req, res)
	}
}
func write_data(vdiname string) {
	clientsdk.Init(MASTER_IP, MASTER_PORT)
	var vdi, err = clientsdk.Open(vdiname)
	if err != nil {
		common.Error("failed to open: %v", err)
		return
	}

	var data = make([]byte, 256, 256)
	data[0] = 'h'
	data[1] = 'e'
	data[2] = 'l'
	data[3] = 'l'
	data[4] = 'o'

	len, err := clientsdk.Write(vdi, data)
	if err != nil {
		common.Error("failed to write: writed=%v err=%v", len, err)
		return
	}
	common.Debug("writed=%v", len)

	data[0] = 'H'
	data[1] = 'E'
	data[2] = 'L'
	data[3] = 'L'
	data[4] = 'O'
	len, err = clientsdk.Write(vdi, data)
	if err != nil {
		common.Error("failed to write: writed=%v err=%v", len, err)
		return
	}
	common.Debug("writed=%v", len)

	clientsdk.Close(vdi)
}
func main() {
	common.Log_init()
	var cmd = os.Args[1]
	switch cmd {
	case "createvdi":
		var vdiname = os.Args[2]
		var vdisize, _ = strconv.ParseUint(os.Args[3], 10, 64)
		create_vdi(vdiname, vdisize)
	case "get_chunkid_meta":
		get_chunkid_meta(os.Args[2:])
	case "write_data":
		var vdiname = os.Args[2]
		write_data(vdiname)
	case "read_write_data":
		var vdiname = os.Args[2]
		read_write_data(vdiname)
	default:
		panic("unknown cmd")
	}
}
