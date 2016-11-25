package protocol

import (
	"time"
	"errors"
	proto "github.com/golang/protobuf/proto"
	
	"net/rpc"
)

import (
	"bytes"
	"encoding/binary"

	"github.com/donnie4w/go-logger/logger"
)

type ProtocolV1 struct {
	Protocol
	Client       *ClientV1
	GatewayAddr    string
	
	rpcGatewayClient       *rpc.Client
}

func (p *ProtocolV1) Exit() {
	if p.Client != nil {
		p.Client.Exit()
	}
}

func (p *ProtocolV1) ProtoHander(cmdType CmdType, dataBtyes []byte) error {
	var err error

	protocolID := cmdType

	var replyCMD CmdType

	if p.Client == nil ||
		p.Client.session.Appkey == nil {
		if protocolID != CmdType_SESSION {
			logger.Fatal("client not vailed app  sessionid:", protocolID)
			return errors.New("client not vailed app")
		}
	}

	var replyData []byte

	switch protocolID {
	case CmdType_SESSION:
		session := SessionMessage{}
		err = proto.Unmarshal(dataBtyes, &session)

		if err != nil {
			logger.Error("ProtoHander: CmdType_SESSION",  err)
			return nil
		}
		
		p.Client.session = session
		
		p.Client.ClientStartTime = time.Now()
		
		logger.Debug("Client Session:", session)
		
		replyCMD = CmdType_SESSIONREPLY
		
		var probability int32
		probability = 100
		
		var retValue bool
		retValue = true
		
		var messageString string
//		messageString = "{\"metric\":\"0\"}"
		messageString=""
		
//		if !p.connectGateway((*session.Appkey)) {
//			probability = 0
//			retValue = true // gateway 不存在 也可以上报
//			messageString = "appkey invailed"
//		} else {			
//			p.Client.session = session
					// 因为目前不需要收集metric数据，所以不再返回metric的概率值
//			enableMetric := p.VailedMeticConfig((*session.Appkey))
	// enableMetric := p.VailedMeticConfig((*session.Appkey))
//			if (enableMetric) {
//				messageString = "{\"metric\":\"1\"}"
//			} else {
//				messageString=""
//			}
//		}
		
		
		reply := &MessageReply{
			Ret:         proto.Bool(retValue),
			Probability: proto.Int32(probability),
			Message:     proto.String(messageString),
		}

		replyData, err = proto.Marshal(reply)

	case CmdType_HTTP_DATA:
		httpData := &HttpData{}
		err = proto.Unmarshal(dataBtyes, httpData)

		if err != nil {
			logger.Error("ProtoHander: CmdType_HTTP_DATA",  err)
			return err
		}

		logger.Debug(httpData)
		p.Client.archiveHttpData(httpData)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)
		
	case CmdType_WEBVIEW_DATA:
		webData := &WebViewData{}
		err = proto.Unmarshal(dataBtyes, webData)

		if err != nil {
			logger.Error("ProtoHander: CmdType_WEBVIEW_DATA",  err)
			return nil
		}

		p.Client.archiveWebViewData(webData)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)


	case CmdType_BATTERY_DATA:
		batteryData := &BatteryData{}
		err = proto.Unmarshal(dataBtyes, batteryData)

		if err != nil {
			logger.Error("ProtoHander: CmdType_BATTERY_DATA",  err)
			return nil
		}

		logger.Debug(p.Client.String(), batteryData)
		p.Client.archiveBatteryData(batteryData)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)	
		
		
	case CmdType_FPS_DATA:
		fallthrough
		
	case CmdType_FPS_CALAMITY_DATA:
		fpsData := &FpsData{}
		err = proto.Unmarshal(dataBtyes, fpsData)

		if err != nil {
			logger.Error("ProtoHander: CmdType_FPS_DATA",  err)
			return nil
		}

		logger.Debug(p.Client.String(), fpsData)
		p.Client.archiveFpsData(fpsData, cmdType)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)	
		
		
	case CmdType_SOCKET_CONNECT:
		connect := &SocketConnect{}
		err = proto.Unmarshal(dataBtyes, connect)

		if err != nil {
			logger.Error("ProtoHander: CmdType_SOCKET_CONNECT",  err)
			return nil
		}
		
		logger.Debug(p.Client.String(), connect)
		p.Client.archiveConnectData(connect)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)

	case CmdType_SOCKET_SENDRECV_DATA:
		sendRecv := &SocketSendRecvData{}
		err = proto.Unmarshal(dataBtyes, sendRecv)

		if err != nil {
			logger.Error("ProtoHander: CmdType_SOCKET_SENDRECV_DATA",  err)
			return nil
		}

		p.Client.archiveSendrecvData(sendRecv)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)

	case CmdType_METHOD_DATA:
		methodData := &MethodData{}
		err = proto.Unmarshal(dataBtyes, methodData)

		if err != nil {
			logger.Error("ProtoHander: CmdType_METHOD_DATA",  err)
			return nil
		}

		logger.Debug(p.Client.String(), methodData)
		p.Client.archiveMethodData(methodData)

		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)	

	case CmdType_UPDATE_NETWORK_STATUS:
		network := AppNetWorkStatus{}
		err = proto.Unmarshal(dataBtyes, &network)
	
		if err != nil {
			logger.Error("ProtoHander: CmdType_UPDATE_NETWORK_STATUS",  err)
			return nil
		}
		
		p.Client.network = network

		logger.Debug(p.Client.String(), network)
		replyCMD = CmdType_RECVDATAREPLY
		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)
		
	case CmdType_METRIC_DATA:
		metricData := &MetricData{}
		err = proto.Unmarshal(dataBtyes, metricData)
		
		if err != nil {
			logger.Error("ProtoHander: CmdType_METRIC_DATA",  err)
			return nil
		}
		
		logger.Debug(p.Client.String(), metricData)
		p.Client.archiveMetricData(metricData)

		reply := &RecvDataReply{
			Recved:  proto.Bool(true),
			Message: proto.String("ok"),
		}

		replyData, err = proto.Marshal(reply)
			
	default:
		logger.Info("default return, not support CMDID", protocolID)
		return err

	}

	if err == nil {

//		logger.Info("reply to client")

		var reply_data_size int32
		reply_data_size = (int32)(len(replyData))

		var reply_session_data []byte

		size_buf := bytes.NewBuffer([]byte{})
		binary.Write(size_buf, binary.BigEndian, reply_data_size)
		reply_session_data = append(reply_session_data, size_buf.Bytes()...)

		var protocolVersion int16
		protocolVersion = 1
		version_buf := bytes.NewBuffer([]byte{})
		binary.Write(version_buf, binary.BigEndian, protocolVersion)
		reply_session_data = append(reply_session_data, version_buf.Bytes()...)

		protocolID_buf := bytes.NewBuffer([]byte{})
		binary.Write(protocolID_buf, binary.BigEndian, replyCMD)
		reply_session_data = append(reply_session_data, protocolID_buf.Bytes()...)

		reply_session_data = append(reply_session_data, replyData...)

		p.Client.Write(reply_session_data)

	} else {
		logger.Info("reply to client err:", err)
	}

	return err
}

func (p *ProtocolV1) connectGateway(appkey string) bool {

	// 目前gateway先不启用，直接返回成功
	return true
	if appkey == "" {
		logger.Warn("appkey is nil ")
		return false;
	}

	var err error
	p.rpcGatewayClient, err = rpc.DialHTTP("tcp", p.GatewayAddr)
	if err != nil {
		logger.Error("connect gateway rpc server failed:", err)
		return false
	}
	var reply bool
	err = p.rpcGatewayClient.Call("Gateway.ConnectGatewayAndVerifyApp", appkey, &reply)
	if err != nil {
		logger.Error("call rep server ConnectGatewayAndVerifyApp failed", err)
		return false
	}

	logger.Info("ConnectGatewayAndVerifyApp ok", reply)
	return reply
}

func (p *ProtocolV1) VailedMeticConfig(appkey string) bool {

	return false;
	
	if p.rpcGatewayClient == nil {
		return false
	}

	var reply bool
	err := p.rpcGatewayClient.Call("Gateway.VerifyMetric", appkey, &reply)
	if err != nil {
		logger.Fatal("VerifyMetric call rep server failed", err)
		return reply
	}
	return reply
}









