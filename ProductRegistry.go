package LifxBridge

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/szammyboi/BitExport"
)

// this should be a part of the light struct

type LifxFeatures struct {
	name      string
	multizone bool
	color     bool
}

func (light *LifxLight) GetProduct() {
	jsonFile, _ := os.Open("./products.json")
	byteValue, _ := io.ReadAll(jsonFile)

	var i []map[string]interface{}
	err := json.Unmarshal(byteValue, &i)
	if err != nil {
		log.Fatal(err)
	}
	products := i[0]["products"].([]interface{})

	header := Header{
		size:        36,
		protocol:    1024,
		addressable: 1,
		tagged:      0,
		source:      2,
		msg_type:    GetVersionPkt,
	}

	var respHeader Header
	var respVersion VersionState
	resp := SendUDP(light.ip, BitExport.ToBytes(header))
	BitExport.FromBytes(resp.data[:36], &respHeader)
	BitExport.FromBytes(resp.data[36:], &respVersion)

	var features LifxFeatures

	for _, product := range products {
		p := product.(map[string]interface{})
		name := p["name"].(string)
		feats := p["features"].(map[string]interface{})
		pid := uint32(p["pid"].(float64))
		if pid == respVersion.product {
			features.multizone = feats["multizone"].(bool)
			features.color = feats["color"].(bool)
			features.name = name
			break
		}
	}
	light.features = features
}
