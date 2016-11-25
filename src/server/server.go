package server

import (
	"bytes"
	"encoding/binary"
	//	"io"
	"encoding/json"
	"net/http"
	"strings"
	"util"

	"abime_lookup"
	"archivedata"
//	"config"
	"fmt"
	"net"
	"os"
	"protocol"
	"time"
	"sync"

	"github.com/donnie4w/go-logger/logger"	
)

/*

1292339 24
1302270 26.3
1334472 25.4

defer close(liveTimeChan)

// https://blog.golang.org/pipelines

*/
import "io/ioutil"

const (
	MAX_CONN_NUM 					= 1024 * 500
	ZOMBIE_CLIENT_NOSEND_TIME  		= time.Minute * 10
)

type ConnSet map[net.Conn]struct{}

type Server struct {
	waitGroup util.WaitGroupWrapper

	tcpListener net.Listener
	svrStatus   SvrStatus
	opts        *abimeOption
	exitChan    chan int
	conns           ConnSet
	mutexConns      sync.Mutex
	
	archive        *archivedata.ArchiveRedis
	
	msgCountChan 	 chan int
}

func NewServer(options *abimeOption) *Server {
	s := &Server{}
	s.opts = options

	s.svrStatus.CurrentClientCount = 0
	s.svrStatus.RunTime = time.Now().Format("2006-01-02 15:04:05")
	s.svrStatus.Name = options.AbimeName	
		
	redisOptions := archivedata.NewRedisOption()
	if s.opts.RedisAddr != "" {
		redisOptions.RedisAddr = s.opts.RedisAddr
	}

	s.archive = archivedata.NewArchiveRidesEx(redisOptions)	
	s.msgCountChan = make(chan int)
	
	return s
}

func (s *Server) Exit() {
	logger.Info("Close TCP Listener begin")
	s.exitChan <- 1
	
	if s.tcpListener != nil {
		logger.Info("Close TCP Listener")
		s.tcpListener.Close()
	}
	
	// shutdown all client
	s.mutexConns.Lock()
	
	for conn, _ := range s.conns {
		conn.Close()
	}
	
	s.conns = nil
		
	s.archive.Exit()
	s.mutexConns.Unlock()
	
	// close(s.msgCountChan)
	// close(s.exitChan)
	
	s.waitGroup.Wait()
}

func (s *Server) serverStatus(w http.ResponseWriter, r *http.Request) {

	if b, err := json.Marshal(s.svrStatus); err == nil {
		var ret string
		ret = string(b)
		w.Write([]byte(ret))
		return
	}

	w.Write([]byte("has error"))
}

func (s *Server) Main() {
	
	go func() {
		PProfInstance(s.opts.Profile)
	}()
	
	//	tcpListener, err := net.Listen("tcp", config.SERVERADDR)
	tcpListener, err := net.Listen("tcp", s.opts.TcpAddr)

	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}

	s.conns = make(ConnSet)
	
	go func() {
		for dot := range s.msgCountChan {
			s.svrStatus.MsgCount += uint64(dot)
		}
	}()
	
	// register to lookup
	// s.regisiterAndReportTolookup()

	go func() {
		for _ = range time.Tick(3 * time.Second) {
			s.regisiterAndReportTolookup()
		}
	}()

	s.tcpListener = tcpListener

	conn_chan := make(chan net.Conn)
	ch_conn_change := make(chan int64)

	go func() {
		for conn_change := range ch_conn_change {
			s.svrStatus.CurrentClientCount += conn_change
		}
	}()
	
//	currentCreatedCount := 0
//	for {
//		if currentCreatedCount >= MAX_CONN_NUM {
//			break;
//		}
//	}

	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range conn_chan {
				ch_conn_change <- 1
				
				s.mutexConns.Lock()
				s.conns[conn] = struct{}{}
				s.mutexConns.Unlock()
							
				s.clientHander(conn)
				
				s.mutexConns.Lock()
				delete(s.conns, conn)
				s.mutexConns.Unlock()
				
				ch_conn_change <- -1
			}
		}()
	}

	s.exitChan = make(chan int, 1)

	s.waitGroup.Wrap(func() {

		go func() {
			for {
				conn, err := s.tcpListener.Accept()
				if err != nil {
					logger.Fatal("Error accept:", err.Error())
					return
				}
				conn_chan <- conn
			}
		}()

		for {
			select {
			case <-s.exitChan:
				logger.Info("exit abime server")
				return
			}
		}
	})

	go func() {
		http.HandleFunc("/status", s.serverStatus)
		http.ListenAndServe(s.opts.HttpAddr, nil)		
	}()
}

func (s * Server) shutdownClentIfVeryIdle(conn net.Conn, liveTimeChan chan time.Time) {
	
	var lastLiveTime time.Time
	lastLiveTime = time.Now()
	
//	exit := false
	
	go func() {
		// check timeout
		for liveTime := range liveTimeChan {
			lastLiveTime = liveTime
//			if exit {
//				return
//			}
		}		
	}()
		
	//  每10s将超过15分钟未发送数据的client踢掉
	go func() {
		for _ = range time.Tick(60 * time.Second) {
			time_now := time.Now()
			
			durationt := time_now.Sub(lastLiveTime)
			if lastLiveTime.Before(time_now) &&
				durationt > ZOMBIE_CLIENT_NOSEND_TIME {
				logger.Debug("Initiative shutdown client ", conn.RemoteAddr().String())
				conn.Close()
//				exit = true
				return
			}		
		}
	}()
}

func (s *Server) clientHander(conn net.Conn) {
	
	/*
		可以考虑添加每个客户端每单位时间(如:15分钟)或者一定时长内没发过数据的可主动踢掉客户端
		1、以让其它的客户端可以链接上来
		2、sdk会有重连机制
		3、减少只链接不发数据的空链接
	*/
	
	defer func() {
//		logger.Info(conn.RemoteAddr().String(), " close socket")
		conn.Close()	
	}()
	
	liveLineChan := make(chan time.Time, 1)
	
	defer close(liveLineChan)
	
	s.shutdownClentIfVeryIdle(conn, liveLineChan)
	
	s.clientMessageLoopHander(conn, liveLineChan)	
}

func (s *Server) clientMessageLoopHander(conn net.Conn, liveLineChan chan time.Time) {
	
	var protoHander *protocol.ProtocolV1

	defer func() {
		if protoHander != nil {
			protoHander.Exit()
		}
	}()
		
	buf := make([]byte, 1024)
	var err error
	var len_data_buf int

	var buf_copy []byte

	var currentProtocolSize int32
	var currentProtocolVersion int16
	var currentProtocolID protocol.CmdType

	var currentProtocolData []byte

	for {
		len_data_buf, err = conn.Read(buf)
		if err != nil {
//			logger.Info("Read Conn Error:", err.Error(), conn.RemoteAddr().String())
			return
		}
		
		liveLineChan <- time.Now()

		buf_copy = append(buf_copy, bytes.NewBuffer(buf[:len_data_buf]).Bytes()...)

		for {
			if currentProtocolSize == 0 {
				if len(buf_copy) < 4 {
					logger.Debug("client protocol size error", currentProtocolSize, "< 4")
					break
				} else {
					var protocolSizeBuf *bytes.Buffer
					protocolSizeBuf = bytes.NewBuffer(buf_copy[:4])
					binary.Read(protocolSizeBuf, binary.BigEndian, &currentProtocolSize)

					if currentProtocolSize <= 0 {
						return
					}
					buf_copy = buf_copy[4:]
				}
			}

			if currentProtocolVersion == 0 {
				if len(buf_copy) < 2 {
					logger.Debug("client protocol version error", currentProtocolVersion, "< 2")
					break
				} else {
					var protocolVersionBuf *bytes.Buffer
					protocolVersionBuf = bytes.NewBuffer(buf_copy[:2])
					binary.Read(protocolVersionBuf, binary.BigEndian, &currentProtocolVersion)

					if currentProtocolVersion != 1 {
						return
					}

					buf_copy = buf_copy[2:]
				}
			}

			if currentProtocolID == 0 {
				if len(buf_copy) < 4 {
//					logger.Debug("client protocol id error", currentProtocolID, "< 4")
					break
				} else {
					var protocolIDBuf *bytes.Buffer
					protocolIDBuf = bytes.NewBuffer(buf_copy[:4])
					binary.Read(protocolIDBuf, binary.BigEndian, &currentProtocolID)

					if currentProtocolID < 0 {
						return
					}

					buf_copy = buf_copy[4:]
				}
			}

			if len(currentProtocolData) == 0 {
				var l int32
				l = (int32)(len(buf_copy))

				if l < currentProtocolSize {
//					logger.Debug("client protocol data error", l, "<", currentProtocolSize)
					break
				} else {

//					logger.Fatal("buf_copy len:", l, "currentProtocolSize:", currentProtocolSize, "buf_size:", len(buf_copy))
					currentProtocolData = buf_copy[:currentProtocolSize]

					if protoHander == nil {
						switch currentProtocolVersion {
						case 1:
							protoHander = &protocol.ProtocolV1{}

							if s.opts.GatewayAddr != "" {
								protoHander.GatewayAddr = s.opts.GatewayAddr
							}
							
							protoHander.Client = protocol.NewClientV1(conn, s.archive)
							
						default:
							logger.Fatal("ERROR: client(%s) bad protocol magic '%d'",
								conn.RemoteAddr(), currentProtocolVersion)
							return
						}
					}

					// hander protocol data
					logger.Debug("PROTOCOL ID", currentProtocolID, "buf:", len(currentProtocolData))

					s.msgCountChan <- 1
					err = protoHander.ProtoHander(currentProtocolID, currentProtocolData)

					if err != nil {
						logger.Fatal("protoHander.ProtoHander", err)
						return
					}

					buf_copy = buf_copy[currentProtocolSize:]
					currentProtocolSize = 0
					currentProtocolVersion = 0
					currentProtocolID = 0
					currentProtocolData = currentProtocolData[0:0]
				}
			}
		}
	}
}

func (s *Server) regisiterAndReportTolookup() {

	registerURL := fmt.Sprintf("http://%s/reportbime", s.opts.LookupAddr)

	client := &http.Client{}

	var abime abime_lookup.Abime
	
	abime.AbimeName 	= s.opts.AbimeName
	abime.ClientCount 	= s.svrStatus.CurrentClientCount
	abime.Addr 			= s.opts.TcpAddr
	abime.MessageCount 	= s.svrStatus.MsgCount
	abime.RunTime 		= s.svrStatus.RunTime

	if b, err := json.Marshal(abime); err == nil {
		body := ioutil.NopCloser(strings.NewReader(string(b)))

		req, _ := http.NewRequest("POST", registerURL, body)

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			//			fmt.Println(string(data), err)

			if string(data) == "ok" {
				//				fmt.Println("report OK")
			}
		} else {
			logger.Error(err)	
		}
	}
}
