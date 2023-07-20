package LifxBridge

import (
	"fmt"
	"math"

	"github.com/szammyboi/BitExport"
)

// hue and lifx can share same api??

// should this pull from server or nah idk
type Light interface {
	UpdateDetails() bool
}

type LifxLight struct {
	ip string
}

func (light LifxLight) UpdateDetails() bool {
	header := Header{
		size:        36,
		protocol:    1024,
		addressable: 1,
		tagged:      0,
		source:      2,
		msg_type:    GetColorPkt,
	}

	var respHeader Header
	var respState LightColorState
	resp := SendUDP(light.ip, BitExport.ToBytes(header))
	BitExport.FromBytes(resp.data[:36], &respHeader)
	BitExport.FromBytes(resp.data[36:], &respState)

	// have bitexport be able to marshal into multiple interfaces (variadic)
	fmt.Printf("Label: %s\n", string(respState.label[:]))
	fmt.Printf(" Power: %d\n", map_range(int64(respState.power), 0, math.MaxUint16, 0, 1))
	fmt.Printf(" Hue: %d\n", map_range(int64(respState.hue), 0, math.MaxUint16, 0, 360))
	fmt.Printf(" Saturation: %.2f\n", map_float(float64(respState.saturation), 0, math.MaxUint16, 0, 1))
	fmt.Printf(" Brightness: %.2f\n", map_float(float64(respState.brightness), 0, math.MaxUint16, 0, 1))

	// kelvin has varying ranges based on device
	// products.json

	return true
}
