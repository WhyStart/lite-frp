package common

import (
	"encoding/binary"
	"encoding/json"
	"lite-frp/tools/log"
)

type ServerData struct {
	ConnId uint64 `json:"connId"`
	Data   []byte `json:"data"` // 数据
}

func ServerByteToData(bytes []byte) ServerData {
	var data ServerData
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		log.Error("byte to data error: %v", err)
	}
	return data
}

func ServerDataToByte(data ServerData) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error("data to byte error: %v", err)
	}
	return bytes
}

func ConvertToByte(data ServerData) []byte {
	requestId := make([]byte, 8)
	binary.LittleEndian.PutUint64(requestId, data.ConnId)
	result := append(requestId, data.Data...)
	return result
}

func ConvertToData(bytes []byte) ServerData {
	connId := binary.LittleEndian.Uint64(bytes[0:8])
	return ServerData{
		ConnId: connId,
		Data:   bytes[8:],
	}
}
