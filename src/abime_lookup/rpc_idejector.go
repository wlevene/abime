package abime_lookup

import (
	//	"encoding/json"
	//	"fmt"
	"net/http"
	"strings"
	//	"time"

	"github.com/donnie4w/go-logger/logger"

	"net"
	"net/rpc"

	//	"config"
)

var id_ejector_pool *IdEjectorPool

type IDEjector int

func (t *IDEjector) ConnectIDEjector(appkey_startID string, reply *bool) error {
	//	*reply = "IDEjector:" + args
	strArray := strings.Split(appkey_startID, ":")
	if len(strArray) != 2 {
		*reply = false
		return nil
	}

	*reply = id_ejector_pool.CreateIdEjector(strArray[0], strArray[1])

	return nil
}

func (t *IDEjector) NewID(appkey string, reply *string) error {

	*reply = id_ejector_pool.NewID(appkey)
	return nil
}

func StartIdEjectorServer() {
	id_ejector_pool = NewIdEjectorPool()

	go func() {
		rpc.Register(new(IDEjector))
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", ":6841")
		if e != nil {
			logger.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
	}()
}
