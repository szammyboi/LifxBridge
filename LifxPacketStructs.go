package LifxBridge

const (
	DiscoveryPkt       uint16 = 2
	DiscoveryRspPkt    uint16 = 3
	GetPowerPkt        uint16 = 20
	GetColorPkt        uint16 = 101
	AcknowledgementPkt uint16 = 45
)

type Header struct {
	size         uint16
	protocol     uint16 `bits:"12"`
	addressable  uint8  `bits:"1"`
	tagged       uint8  `bits:"1"`
	origin       uint8  `bits:"2"`
	source       uint32
	target       [8]uint8
	_            [6]uint8
	res_required uint8 `bits:"1"`
	ack_required uint8 `bits:"1"`
	_            uint8 `bits:"6"`
	sequence     uint8
	_            uint64
	msg_type     uint16
	_            uint16
}

type LightColorState struct {
	hue        uint16
	saturation uint16
	brightness uint16
	kelvin     uint16
	_          [2]byte
	power      uint16
	label      [32]byte
	_          [8]byte
}

func AckReceived(header Header) bool {
	return header.msg_type == AcknowledgementPkt
}
