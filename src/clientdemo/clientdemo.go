package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func main() {

	var protocolSize int32
	protocolSize = 5

	var protocolVersion int16
	protocolVersion = 1

	var protocolID int32
	protocolID = 1

	size_buf := bytes.NewBuffer([]byte{})
	binary.Write(size_buf, binary.BigEndian, protocolSize)
	fmt.Println(size_buf.Bytes())

	version_buf := bytes.NewBuffer([]byte{})
	binary.Write(version_buf, binary.BigEndian, protocolVersion)
	fmt.Println(version_buf.Bytes())

	//	header_buf := size_buf.Bytes() + version_buf.Bytes()
	header_buf := make([]byte, 0)

	header_buf = append(header_buf, size_buf.Bytes()...)
	header_buf = append(header_buf, version_buf.Bytes()...)

	id_buf := bytes.NewBuffer([]byte{})
	binary.Write(id_buf, binary.BigEndian, protocolID)

	header_buf = append(header_buf, id_buf.Bytes()...)

	fmt.Println(len(header_buf), ":", header_buf)

	proto1 := append(header_buf, []byte("hello")...)
	fmt.Println(len(proto1), ":", proto1)
	//	return
	conn, err := net.Dial("tcp", "0.0.0.0:9372")
	checkError(err)
	_, err = conn.Write(proto1)

	_, err = conn.Write([]byte("Hello Server"))
	_, err = conn.Write([]byte("Hello Server"))
	//	_, err = conn.Write([]byte("Hello Server"))
	//	checkError(err)

	time.Sleep(1 * time.Second)

	//	_, err = conn.Write([]byte("hahaha"))
	//	_, err = conn.Write([]byte("hahahaa"))
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
