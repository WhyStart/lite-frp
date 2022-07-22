package server

import (
	"fmt"
	"lite-frp/config"
	"lite-frp/tools/common"
	"lite-frp/tools/crypt"
	"lite-frp/tools/log"
	"net"
	"os"
	"strconv"
	"sync"
)

type Service struct {
	cfg config.ServerCommonConf
}

var clientConn sync.Map
var maxKey = 0

func NewService(config config.ServerCommonConf) *Service {
	return &Service{
		cfg: config,
	}
}

func (svr *Service) Run() {
	log.InitLog()
	address := net.JoinHostPort(svr.cfg.BindAddr, strconv.Itoa(svr.cfg.BindPort))
	listen, err := net.Listen(svr.cfg.Type, address)
	if err != nil {
		log.Error("listen error：%v", err)
		os.Exit(1)
	}
	log.Info("lite_server服务端启动成功：[%s %s:%d]", svr.cfg.Type, svr.cfg.BindAddr, svr.cfg.BindPort)
	for {
		if c, err := listen.Accept(); err != nil {
			log.Error("accept error: %v", err)
		} else {
			log.Info("接收连接")
			content, err := common.GetContent(c)
			if err != nil {
				log.Info("客户端连接关闭,%v", err)
				c.Close()
			}
			data := common.ClientByteToData(content)
			if svr.verifyKey(data) {
				log.Info("校验通过，创建连接")
				c.Write([]byte("1"))
				go svr.connHandle(c, data)
				go clientDataHandle(c)
			} else {
				log.Info("校验未通过，关闭连接")
				c.Write([]byte("connection validation error"))
				c.Close()
			}
		}
	}
}

func (svr *Service) verifyKey(data common.ClientData) bool {
	return crypt.Md5(svr.cfg.Key) == data.Key
}

func (svr *Service) connHandle(conn net.Conn, data common.ClientData) {
	if data.RequestType == common.Create {
		log.Info("创建连接")
		go svr.createConnection(data, conn)
	}
}

func clientDataHandle(conn net.Conn) {
	for {
		content, err := common.GetContent(conn)
		if err != nil {
			log.Error("读取客户端数据错误：%v", err)
			break
		}
		log.Info("读取客户端数据")
		serverData := common.ConvertToData(content)
		load, ok := clientConn.Load(serverData.ConnId)
		if ok {
			log.Info("发送数据到浏览器")
			(load.(net.Conn)).Write(serverData.Data)
		}
	}
}

func (svr *Service) createConnection(data common.ClientData, conn net.Conn) {
	listen, err := net.Listen(svr.cfg.Type, fmt.Sprintf("%s:%d", svr.cfg.BindAddr, data.Port))
	if err != nil {
		log.Error("error listen:%v", err)
		return
	}
	log.Info("创建连接成功，remote port为：%d", data.Port)
	defer listen.Close()
	for {
		if c, err := listen.Accept(); err == nil {
			log.Info("请求地址：%v", c.LocalAddr())

			go func() {
				connId := common.GetConnId()
				clientConn.Store(connId, c)
				for {
					content, err := common.GetContent(c)
					if err != nil {
						log.Error("读取浏览器数据错误：%v", err)
						clientConn.Delete(connId)
						c.Close()
						break
					}
					if len(content) == 0 {
						log.Info("浏览器请求数据读取完毕关闭连接")
						c.Close()
						clientConn.Delete(connId)
						break
					}
					serverData := common.ServerData{Data: content, ConnId: connId}
					bytes := common.ConvertToByte(serverData)
					_, err = conn.Write(bytes)
					if err != nil {
						clientConn.Delete(connId)
						c.Close()
						log.Error("写入到客户端错误")
					}
				}
			}()
		} else {
			log.Error("accept error:%v", err)
		}
	}
}
