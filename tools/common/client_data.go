package common

import (
	"encoding/json"
	"lite-frp/tools/log"
	"net"
)

type ClientData struct {
	RequestType int    `json:"requestType"`
	Key         string `json:"key"`
	Port        int    `json:"port"`
	Data        []byte `json:"data"` // 数据
}

func connToClientData(conn net.Conn) ClientData {
	return ClientData{}
}

func ClientDataToByte(data ClientData) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error("data to byte error: %v", err)
	}
	//requestId := make([]byte, 8)
	//binary.LittleEndian.PutUint64(requestId, uint64(data.RequestId))
	//
	//serviceId := make([]byte, len(data.ConnId))
	//binary.LittleEndian.PutUint64(requestId, uint64(data.RequestId))
	//
	//ret := append(requestId, serviceId...)
	//ret = append(ret, data.Data...)

	return bytes
}

func ClientByteToData(bytes []byte) ClientData {
	var data ClientData
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		log.Error("byte to data error: %v", err)
	}
	return data
}
