package common

import (
	"net"
)

func Generate_chunkserver_id(ips string, port int) uint64 {
	var ip = net.ParseIP(ips)
	var ip_int = ip[3]<<24 | ip[2]<<16 | ip[1]<<8 | ip[0]

	return uint64(ip_int)<<32 | uint64(uint32(port))
}
func Unpack_chunkserver_id(chunkserver_id uint64) (string, int) {
	var ip_int = chunkserver_id >> 32
	var a = ip_int >> 24
	var b = (ip_int >> 16) & 0xff
	var c = (ip_int >> 8) & 0xff
	var d = (ip_int) & 0xff
	var ip = net.IPv4(byte(a), byte(b), byte(c), byte(d)).String()

	var port = int(chunkserver_id & 0xffff)

	return ip, port
}
func Get_vdiid_from_chunkid(chunkid uint64) uint32 {
	return uint32(chunkid >> 32)
}
func Generate_chunkid(vdiid uint32, idx uint32) uint64 {
	return uint64(vdiid)<<32 | uint64(idx)
}
