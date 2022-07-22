package common

import (
	"io"
	"net"
	"sync"
)

var (
	mu     sync.Mutex
	connId uint64
)

func GetContent(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4104) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return []byte{}, err
			}
			break
		}
		buf = append(buf, tmp[:n]...)
		if n < len(tmp) || len(buf) == 4104 {
			break
		}
	}
	return buf, nil
}

func GetConnId() uint64 {
	mu.Lock()
	defer mu.Unlock()
	if connId == ^uint64(0) {
		connId = 0
	}
	connId += 1
	return connId
}
