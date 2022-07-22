package client

import (
	"encoding/json"
	"lite-frp/config"
	"lite-frp/tools/common"
	"lite-frp/tools/crypt"
	"lite-frp/tools/log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	cfg        config.ClientCommonConf
	serverConn net.Conn
}

func NewService(config config.ClientCommonConf) *Service {
	return &Service{
		cfg: config,
	}
}

var con sync.Map

func (svr *Service) Run() {
	log.InitLog()
	for {
		address := net.JoinHostPort(svr.cfg.ServerAddr, strconv.Itoa(svr.cfg.ServerPort))
		conn, err := net.Dial(svr.cfg.Type, address)
		if err != nil {
			log.Warn("与Server连接失败：dial %s %s:%d，10秒后重试", svr.cfg.Type, svr.cfg.ServerAddr, svr.cfg.ServerPort)
			time.Sleep(10 * time.Second)
			continue
		}
		svr.serverConn = conn
		log.Info("向Server发送校验连接数据")
		svr.sendDataToServer(common.ClientData{
			RequestType: common.Create,
			Key:         crypt.Md5(svr.cfg.SK),
			Port:        svr.cfg.RemotePort,
		})
		buf := make([]byte, 1)
		read, err := conn.Read(buf)
		if read != 1 || err != nil {
			log.Error("读取Server返回数据错误：%v", err)
		}
		if string(buf[0]) == "1" {
			log.Info("校验通过，连接Server成功")
			svr.serverDataHandle()
		} else {
			log.Error("校验失败，关闭连接")
			conn.Close()
			//os.Exit(1)
		}
	}
}

var mu sync.Mutex

func (svr *Service) sendDataToServer(data common.ClientData) {
	log.Info("发送数据到服务器")
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error("数据转换json失败：%v", err)
	}
	svr.serverConn.Write(bytes)
}

var socket sync.Map

func (svr *Service) serverDataHandle() {
	for {
		content, err := common.GetContent(svr.serverConn)
		if err != nil {
			log.Error("读取Server连接数据错误：%v", err)
			svr.serverConn.Close()
			break
		}
		serverData := common.ConvertToData(content)
		poolConn, ok := socket.Load(serverData.ConnId)
		if !ok {
			log.Info("创建本地socket连接")
			poolConn, err = svr.createLocalSocket(serverData.ConnId)
			if err != nil {
				log.Error("创建本地连接失败：%v", err)
				return
			}
		}
		log.Info("向本地socket连接发送Server返回的数据")
		(poolConn.(net.Conn)).Write(serverData.Data)
	}
}

func (svr *Service) createLocalSocket(connId uint64) (net.Conn, error) {
	address := net.JoinHostPort(svr.cfg.LocalIP, strconv.Itoa(svr.cfg.LocalPort))
	conn, err := net.Dial(svr.cfg.Type, address)
	if err != nil {
		return nil, err
	}
	socket.Store(connId, conn)
	go func() {
		for {
			localSocketContent, err := common.GetContent(conn)
			if err != nil {
				conn.Close()
				socket.Delete(connId)
				log.Error("读取本地连接数据错误，%v", err)
				break
			}
			if len(localSocketContent) == 0 {
				log.Info("本地socket数据读取完毕关闭连接")
				conn.Close()
				socket.Delete(connId)
				break
			}
			log.Info("读取本地socket数据")
			func() {
				mu.Lock()
				defer mu.Unlock()
				log.Info("本地socket数据发送到Server")
				serverData := common.ServerData{Data: localSocketContent, ConnId: connId}
				bytes := common.ConvertToByte(serverData)
				if err != nil {
					log.Error("发送数据到Server错误：%v", err)
				}
				_, err = svr.serverConn.Write(bytes)
				if err != nil {
					log.Error("发送数据到Server错误：%v", err)
					conn.Close()
					socket.Delete(connId)
				}
			}()
		}
	}()
	return conn, nil
}
