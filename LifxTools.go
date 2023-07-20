package LifxBridge

import (
	"log"
	"net"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type NetworkResponse struct {
	addr string
	data []byte
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// use interfaces
// always sending and receiving from the same port means I wouldn't have to have a timeout
// discovery is the only function that needs to wait for mutliple responses though

func MakeUDPConn(dest_ip string, data []byte) *net.PacketConn {
	conn, err := net.ListenPacket("udp", ":0")
	CheckErr(err)

	dest, err := net.ResolveUDPAddr("udp", dest_ip)
	CheckErr(err)

	//fmt.Printf("Sent %d bytes from %s\n", len(data), conn.LocalAddr().String())
	n, err := conn.WriteTo(data, dest)
	CheckErr(err)
	if n != len(data) {
		log.Fatal("All bytes were not written!")
	}

	return &conn
}

// parallelization
// currently singular
// maybe have a thread always running to send messages
// then just change destination
func SendUDP(dest_ip string, data []byte) NetworkResponse {
	conn := *MakeUDPConn(dest_ip, data)

	buffer := make([]byte, 4096)
	n, addr, err := conn.ReadFrom(buffer)
	CheckErr(err)
	//fmt.Printf("Received %d bytes from %s\n", n, addr.String())

	resp := NetworkResponse{
		addr: addr.String(),
		data: buffer[:(n + 1)],
	}

	conn.Close()
	return resp
}

func SendUDPMulti(dest_ip string, data []byte) []NetworkResponse {
	addrs := mapset.NewSet[string]()
	resps := make([]NetworkResponse, 0)
	conn := *MakeUDPConn(dest_ip, data)

	deadline := time.Now().Add(time.Millisecond * 500)
	conn.SetReadDeadline(deadline)

	for {
		buffer := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(buffer)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			break
		}
		CheckErr(err)

		if addrs.Contains(addr.String()) {
			continue
		} else {
			addrs.Add(addr.String())
		}

		//fmt.Printf("Received %d bytes from %s\n", n, addr.String())

		resp := NetworkResponse{
			addr: addr.String(),
			data: buffer[:(n + 1)],
		}

		resps = append(resps, resp)

	}

	conn.Close()
	return resps
}

func map_range(x int64, in_min int64, in_max int64, out_min int64, out_max int64) int64 {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}

func map_float(x float64, in_min float64, in_max float64, out_min float64, out_max float64) float64 {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}

/*

 */
