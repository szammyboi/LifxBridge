package LifxBridge

import (
	"fmt"

	"github.com/szammyboi/BitExport"
)

func Discovery() {
	header := Header{
		size:        36,
		protocol:    1024,
		addressable: 1,
		tagged:      1,
		source:      2,
		msg_type:    DiscoveryPkt,
	}

	resps := SendUDPMulti(LifxBroadcast, BitExport.ToBytes(header))

	// update lights in database
	// just existence
	// server code will not see any of this
	// nah it will run this so idk
	for _, resp := range resps {
		// new light func to auto get features
		light := LifxLight{
			ip: resp.addr,
		}
		light.features = light.GetProduct()
		fmt.Println(light)
		light.UpdateDetails()
	}
}
