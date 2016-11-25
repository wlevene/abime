package main

import (
	"fmt"
	"net/rpc"
	"time"
	"sync"  
//	"os"
//	"os/signal"
//	"syscall"
)


var once sync.Once
var Gid int

func setup() {
	Gid++
	fmt.Println("Called once")
}

func doprint() {
	once.Do(setup)
	fmt.Println("doprint()...")
}


func main() {
	
	
	now := time.Now().UnixNano()
			
 	fmt.Println(fmt.Sprintf("%d", now))
	
	
//	i, _ = syscall.ForkExec(os.Args[0], os.Args, nil)
//	j , _= syscall.ForkExec(os.Args[0], os.Args, nil)
	
	
//	time.Sleep(1 *  time.Second)
	
//	fmt.Println("i:", i, "  j:", j)
	return
//	listenerFile, err := listener.File()
//	if err != nil {
//    	fmt.Println("Fail to get socket file descriptor:", err)
//	}
//	listenerFd := listenerFile.Fd()

//	// Set a flag for the new process start
//	processos.Setenv("_GRACEFUL_RESTART", "true")
	
//	execSpec := &syscall.ProcAttr{
//    	Env:   os.Environ(),
//    	Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
//	}
	
//	// Fork exec the new version of your
//	serverfork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	
//	syscall.ForkExec()
	
	return
	
	t1 := time.Now()
	
	time.Sleep(time.Second * 2)
	
	t2 := time.Now().Sub(t1)
	
	
	
	fmt.Println(int64(t2.Seconds()))
	
	
	return;
	
	for i := 0; i<15; i++  {
		client1, err1 := rpc.DialHTTP("tcp", ":6848")
		if err1 != nil {
			fmt.Println("链接rpc服务器失败:", err1)
			return
		}
		var reply1 bool
		err1 = client1.Call("Gateway.ConnectGatewayAndVerifyApp", "1212", &reply1)

		if err1 != nil {
			fmt.Println("调用远程服务失败", err1)
			return
		}
		fmt.Println(reply1)
	
	
		err1 = client1.Call("Gateway.VerifyMetric", "1212", &reply1)
		if err1 != nil {
			fmt.Println("---调用远程服务失败", err1)
			return
		}
		fmt.Println("--",reply1)
	}
	
	
	
	return;
	
	go doprint()
	go doprint()
	go doprint()
	go doprint()

	time.Sleep(time.Second)
	fmt.Println("Gid:", Gid)
	
	
	return

	for i := 0; i < 100; i++ {
		t := time.Now().UnixNano()
		fmt.Println(t)
	}

	return
	a := []int{1, 2, 3, 4}

	index := 2
	abime := 6

	copyPool := a
	fmt.Println("1:", copyPool)
	a = append(copyPool[:index], abime)
	fmt.Println("2:", a, copyPool)
	a = append(a, copyPool[index:]...)
	fmt.Println("3:", a)
	copyPool = copyPool[0:0]

	return

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("链接rpc服务器失败:", err)
		return
	}
	var reply string
	fmt.Println(time.Now())
	err = client.Call("Watcher.GetInfo", "1", &reply)
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("调用远程服务失败", err)
		return
	}
	fmt.Println("远程服务返回结果：", reply)
	err = client.Call("Watcher.GetInfo", "1", &reply)
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("调用远程服务失败", err)
		return
	}
	fmt.Println("远程服务返回结果：", reply)
	err = client.Call("Watcher.GetInfo", "1", &reply)
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("调用远程服务失败", err)
		return
	}
	fmt.Println("远程服务返回结果：", reply)
}
