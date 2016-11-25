package util

import (
	"net"
)

//func getWWWIPAddr() string {
//	conn, err := net.Dial("udp", "baidu.com:80")
//	if err != nil {
//		//		fmt.Println(err.Error())
//		return ""
//	}
//	defer conn.Close()
//	return strings.Split(conn.LocalAddr().String(), ":")[0]
//}

func getIPAddr() string {

	var ipaddr string
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ipaddr
	}

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipaddr = ipnet.IP.String()
				return ipaddr
			}

		}
	}

	return ipaddr
}

