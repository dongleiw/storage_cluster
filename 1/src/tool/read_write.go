package main

import (
	"clientsdk"
	"common"
	//"fmt"
)

func read_write_data(vdiname string) {
	clientsdk.Init(MASTER_IP, MASTER_PORT)
	var vdi, err = clientsdk.Open(vdiname)
	if err != nil {
		common.Error("failed to open: %v", err)
		return
	}

	if false {
		var data = make([]byte, 256, 256)
		// write
		data[0] = 'h'
		data[1] = 'e'
		data[2] = 'l'
		data[3] = 'l'
		data[4] = 'o'

		writelen, err := clientsdk.Write(vdi, data)
		if err != nil {
			common.Error("failed to write: writed=%v err=%v", writelen, err)
			return
		}
		common.Debug("writed=%v", writelen)

		// write
		data[0] = 'H'
		data[1] = 'E'
		data[2] = 'L'
		data[3] = 'L'
		data[4] = 'O'
		writelen, err = clientsdk.Write(vdi, data)
		if err != nil {
			common.Error("failed to write: writed=%v err=%v", writelen, err)
			return
		}
		common.Debug("writed=%v", writelen)
	}

	// read
	clientsdk.Seek(vdi, 0, clientsdk.SEEK_SET)
	{
		var data, err = clientsdk.Read(vdi, 256)
		if err != nil {
			common.Error("failed to read: %v", err)
			return
		}
		common.Debug("read [%v]", string(data))
	}
	{
		var data, err = clientsdk.Read(vdi, 256)
		if err != nil {
			common.Error("failed to read: %v", err)
			return
		}
		common.Debug("read [%v]", string(data))
	}

	clientsdk.Close(vdi)
}
