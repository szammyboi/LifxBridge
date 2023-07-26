package LifxBridge

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/szammyboi/BitExport"
)

// this should be a part of the light struct

type LifxFeatures struct {
	Name      string
	MultiZone bool
}

func (light LifxLight) GetProduct() LifxFeatures {
	// figure out how to not do this
	jsonFile, _ := os.Open("/home/szammy/code/go/Eos/Tools/LifxBridge/products.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var i []map[string]interface{}
	json.Unmarshal(byteValue, &i)
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
			features.MultiZone = feats["multizone"].(bool)
			features.Name = name
			break
		}
	}

	return features
}
